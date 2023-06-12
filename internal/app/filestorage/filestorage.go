package filestorage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/kripsy/shortener/internal/app/models"
	"github.com/kripsy/shortener/internal/app/utils"
	"go.uber.org/zap"
)

type FileStorage struct {
	memoryStorage map[string]string
	fileName      string
	myLogger      *zap.Logger
}

func InitFileStorageFile(fileName string, myLogger *zap.Logger) (*FileStorage, error) {
	if fileName == "" {
		return nil, errors.New("fileName is empty")
	}
	memoryStorage := map[string]string{}

	fs := &FileStorage{
		memoryStorage,
		fileName,
		myLogger,
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
	for _, event := range events {
		fs.memoryStorage[event.OriginalURL] = event.ShortURL
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

func (fs *FileStorage) CreateOrGetFromStorage(ctx context.Context, url string) (string, error) {

	for originalURL, shortURL := range fs.memoryStorage {
		if originalURL == url {
			return shortURL, nil
		}
	}

	shortURL, err := utils.CreateShortURL()
	if err != nil {
		return "", err
	}
	event := models.NewEvent(shortURL, url)
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

	fs.memoryStorage[url] = shortURL

	return shortURL, nil
}

func (fs *FileStorage) GetOriginalURLFromStorage(ctx context.Context, shortURL string) (string, error) { //([]models.Event, error)

	var val string
	ok := false
	// for every key from MYMEMORY check our shortURL. If exist set `val = k` and `ok = true`

	for k, v := range fs.memoryStorage {
		if v == string(shortURL) {
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

func (fs *FileStorage) CreateOrGetBatchFromStorage(ctx context.Context, batchURL *models.BatchURL) (*models.BatchURL, error) {
	fs.myLogger.Error("Start CreateOrGetBatchFromStorage")
	for k, v := range *batchURL {
		shortURL, err := fs.CreateOrGetFromStorage(context.Background(), v.OriginalURL)
		if err != nil {
			return nil, err
		}
		(*batchURL)[k].ShortURL = shortURL
		(*batchURL)[k].OriginalURL = ""
	}
	return batchURL, nil
}
