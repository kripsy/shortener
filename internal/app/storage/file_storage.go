package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type FileStorage struct {
	FileName string
}

func InitFileStorageFile(fileName string) *FileStorage {
	fs := &FileStorage{
		FileName: fileName,
	}
	return fs
}

type Event struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
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

func (p *Producer) WriteEvent(event Event) error {
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
		fmt.Println(fileName)
		fmt.Println(err)
		myLogger.Warn("errror create file to consumer")
		return nil, err
	}

	c := &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}

	return c, nil
}

func (c *Consumer) ReadEvent() (*Event, error) {
	event := &Event{}

	if err := c.decoder.Decode(&event); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}

func NewEvent(shortURL string, originalURL string) *Event {
	e := &Event{
		UUID:        int(uuid.New().ID()),
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}
	return e
}

func addURL(events []Event, fileName string, myLogger *zap.Logger) error {

	Producer, err := NewProducer(fileName, myLogger)
	if err != nil {
		myLogger.Error("cannot create producer")
		return err
	}
	defer Producer.Close()
	for _, event := range events {
		if err := Producer.WriteEvent(event); err != nil {
			myLogger.Error("error write event")
			return nil
		}
	}
	return nil
}

func readURL(fileName string, myLogger *zap.Logger) ([]Event, error) {
	Consumer, err := NewConsumer(fileName, myLogger)
	events := []Event{}
	if err != nil {
		myLogger.Error("cannot create Consumer")
		return nil, err
	}

	var event *Event
	for err == nil {
		event, err = Consumer.ReadEvent()
		if err == nil {
			// fmt.Println(event)
			myLogger.Debug("New event", zap.Int("UUID", event.UUID),
				zap.String("short_url", event.ShortURL),
				zap.String("original_url", event.OriginalURL))
			events = append(events, *event)
		}
	}

	return events, nil
}
