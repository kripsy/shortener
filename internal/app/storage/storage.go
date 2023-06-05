package storage

import (
	"fmt"

	"github.com/kripsy/shortener/internal/app/utils"
	"go.uber.org/zap"
)

type Storage struct {
	storage         map[string]string
	MyLogger        *zap.Logger
	fileStorageName string
}

func InitStorageFromFile(storage map[string]string, fs *FileStorage, myLogger *zap.Logger) error {
	events, err := readURL(fs.FileName, myLogger)
	if err != nil {
		myLogger.Warn("error read URLs")
		return err
	}
	for _, v := range events {
		storage[v.OriginalURL] = v.ShortURL
	}
	return nil
}

func InitStorage(initValue map[string]string, fs *FileStorage, myLogger *zap.Logger) *Storage {
	m := Storage{
		storage:         initValue,
		MyLogger:        myLogger,
		fileStorageName: fs.FileName,
	}
	InitStorageFromFile(m.storage, fs, m.MyLogger)
	return &m
}

func (m *Storage) CreateOrGetFromStorage(url string) (string, error) {
	// If the key exists
	val, ok := m.storage[url]
	if !ok {
		// input into our storage
		val, err := utils.CreateShortURL()
		if err != nil {
			return "", err
		}
		m.storage[url] = val

		e := make([]Event, 0)

		if ne := NewEvent(val, url); ne != nil {
			e = append(e, *ne)
		}

		if len(e) > 0 {
			addURL(e, m.fileStorageName, m.MyLogger)
		}

		return val, nil
	}
	return val, nil
}

func (m Storage) GetFromStorage(url string) (string, error) {

	var val string
	ok := false
	// for every key from MYMEMORY check our shortURL. If exist set `val = k` and `ok = true`

	for k, v := range m.storage {
		if v == string(url) {
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
