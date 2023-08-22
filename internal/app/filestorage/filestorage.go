package filestorage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/kripsy/shortener/internal/app/models"
	"github.com/kripsy/shortener/internal/app/utils"
	"go.uber.org/zap"
)

type FileStorage struct {
	memoryStorage map[string]models.Event
	fileName      string
	myLogger      *zap.Logger
	rwmutex       *sync.RWMutex
}

func InitFileStorageFile(fileName string, myLogger *zap.Logger) (*FileStorage, error) {
	if fileName == "" {
		return nil, errors.New("fileName is empty")
	}
	memoryStorage := map[string]models.Event{}
	rwmutex := &sync.RWMutex{}

	fs := &FileStorage{
		memoryStorage,
		fileName,
		myLogger,
		rwmutex,
	}
	err := fs.fillMemoryStorage()
	if err != nil {
		return nil, err
	}
	return fs, nil
}

func (fs *FileStorage) fillMemoryStorage() error {
	events, err := fs.readEventsFromFile()
	if err != nil {
		fs.myLogger.Warn("error fillMemoryStorage", zap.String("msg", err.Error()))
		return err
	}
	fs.rwmutex.Lock()
	defer fs.rwmutex.Unlock()
	for _, event := range events {
		fs.memoryStorage[event.OriginalURL] = event
	}
	return nil
}

type Producer struct {
	file    *os.File
	encoder json.Encoder
}

func NewProducer(fileName string, myLogger *zap.Logger) (*Producer, error) {

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		myLogger.Warn("errror create file to producer")
		fmt.Println(err)
		return nil, err
	}

	p := &Producer{
		file:    file,
		encoder: *json.NewEncoder(file),
	}

	return p, nil
}

func (p *Producer) WriteEvent(event models.Event) error {
	return p.encoder.Encode(event)
}

func (p *Producer) Close() error {
	return p.file.Close()
}

type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewConsumer(fileName string, myLogger *zap.Logger) (*Consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		myLogger.Warn("errror create file to consumer", zap.String("msg", err.Error()))
		return nil, err
	}

	c := &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}

	return c, nil
}

func (c *Consumer) ReadEvent() (*models.Event, error) {
	event := &models.Event{}

	if err := c.decoder.Decode(&event); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}

func (fs *FileStorage) CreateOrGetFromStorage(ctx context.Context, url string, userID int) (string, error) {
	fs.rwmutex.Lock()
	defer fs.rwmutex.Unlock()
	for originalURL, event := range fs.memoryStorage {
		if originalURL == url {
			return event.ShortURL, nil
		}
	}

	// shortURL, err := utils.CreateShortURL()
	shortURL, err := utils.CreateShortURLWithoutFmt()
	if err != nil {
		return "", err
	}
	event := models.NewEvent(shortURL, url, userID)
	Producer, err := NewProducer(fs.fileName, fs.myLogger)
	if err != nil {
		fs.myLogger.Error("cannot create producer")
		return "", err
	}
	defer Producer.Close()

	if err = Producer.WriteEvent(*event); err != nil {
		fs.myLogger.Error("error write event", zap.String("msg", err.Error()))
		return "", err
	}

	fs.memoryStorage[url] = *event

	return shortURL, nil
}

func (fs *FileStorage) GetOriginalURLFromStorage(ctx context.Context, shortURL string) (string, error) { //([]models.Event, error)
	fs.rwmutex.RLock()
	defer fs.rwmutex.RUnlock()
	var val string
	ok := false
	// for every key from MYMEMORY check our shortURL. If exist set `val = k` and `ok = true`

	for k, v := range fs.memoryStorage {
		if v.ShortURL == string(shortURL) {
			ok = true
			val = k
			break
		}
	}
	if !ok {
		// key not exist
		return "", fmt.Errorf("not exists")
	}
	// If the key exists
	return val, nil
}

func (fs *FileStorage) readEventsFromFile() ([]models.Event, error) {
	fs.rwmutex.RLock()
	defer fs.rwmutex.RUnlock()
	Consumer, err := NewConsumer(fs.fileName, fs.myLogger)
	events := []models.Event{}
	if err != nil {
		fs.myLogger.Error("cannot create Consumer")
		return nil, err
	}
	var event *models.Event
	for err == nil {
		event, err = Consumer.ReadEvent()
		if err == nil {
			fs.myLogger.Debug("New event", zap.Int("UUID", event.UUID),
				zap.String("short_url", event.ShortURL),
				zap.String("original_url", event.OriginalURL))
			events = append(events, *event)
		}
	}

	return events, nil
}

func (fs *FileStorage) Close() {

}
func (fs *FileStorage) Ping() error {
	return nil
}

func (fs *FileStorage) CreateOrGetBatchFromStorage(ctx context.Context, batchURL *models.BatchURL, userID int) (*models.BatchURL, error) {
	fs.myLogger.Debug("Start CreateOrGetBatchFromStorage")
	for k, v := range *batchURL {
		shortURL, err := fs.CreateOrGetFromStorage(context.Background(), v.OriginalURL, userID)
		if err != nil {
			return nil, err
		}
		(*batchURL)[k].ShortURL = shortURL
		(*batchURL)[k].OriginalURL = ""
	}
	return batchURL, nil
}

func (fs *FileStorage) GetUserByID(ctx context.Context, ID int) (*models.User, error) {
	events, err := fs.readEventsFromFile()
	if err != nil {
		fs.myLogger.Debug("Error read events", zap.String("msg", err.Error()))
		return nil, err
	}
	for _, v := range events {
		if v.UserID == ID {
			return &models.User{
				ID: v.UserID,
			}, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (fs *FileStorage) RegisterUser(ctx context.Context) (*models.User, error) {
	return &models.User{
		ID: int(uuid.New().ID()),
	}, nil
}

func (fs *FileStorage) GetBatchURLFromStorage(ctx context.Context, userID int) (*models.BatchURL, error) {
	batchURL := &models.BatchURL{}
	events, err := fs.readEventsFromFile()
	if err != nil {
		fs.myLogger.Debug("Error read events", zap.String("msg", err.Error()))
		return nil, err
	}
	for _, v := range events {
		if v.UserID == userID {
			event := &models.Event{
				ShortURL:    v.ShortURL,
				OriginalURL: v.OriginalURL,
			}
			*batchURL = append(*batchURL, *event)
		}
	}
	return batchURL, nil
}

func (fs *FileStorage) DeleteSliceURLFromStorage(ctx context.Context, shortURL []string, userID int) error {
	fmt.Println("Not implemented yet")

	return nil
}
