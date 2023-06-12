package inmemorystorage

import (
	"context"
	"fmt"

	"github.com/kripsy/shortener/internal/app/models"
	"github.com/kripsy/shortener/internal/app/utils"
	"go.uber.org/zap"
)

type InMemoryStorage struct {
	storage  map[string]string
	myLogger *zap.Logger
}

// func InitStorageFromFile(storage map[string]string, fs *FileStorage, myLogger *zap.Logger) error {
// 	events, err := readURL(fs.FileName, myLogger)
// 	if err != nil {
// 		myLogger.Warn("error read URLs")
// 		return err
// 	}
// 	for _, v := range events {
// 		storage[v.OriginalURL] = v.ShortURL
// 	}
// 	return nil
// }

func InitInMemoryStorage(initValue map[string]string, myLogger *zap.Logger) (*InMemoryStorage, error) {
	m := &InMemoryStorage{
		storage:  initValue,
		myLogger: myLogger,
	}
	return m, nil
}

func (m *InMemoryStorage) CreateOrGetFromStorage(ctx context.Context, url string) (string, error) {
	// If the key exists
	val, ok := m.storage[url]
	if !ok {
		// input into our storage
		val, err := utils.CreateShortURL()
		if err != nil {
			return "", err
		}
		m.storage[url] = val

		// e := make([]models.Event, 0)

		// if ne := models.NewEvent(val, url); ne != nil {
		// 	e = append(e, *ne)
		// }

		// if len(e) > 0 {
		// 	addURL(e, m.fileStorageName, m.MyLogger)
		// }

		return val, nil
	}
	return val, nil
}

func (m InMemoryStorage) GetOriginalURLFromStorage(ctx context.Context, shortURL string) (string, error) {

	var val string
	ok := false
	// for every key from MYMEMORY check our shortURL. If exist set `val = k` and `ok = true`

	for k, v := range m.storage {
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

func (m InMemoryStorage) Close() {
}

func (m InMemoryStorage) Ping() error {
	return nil
}

func (m InMemoryStorage) CreateOrGetBatchFromStorage(ctx context.Context, batchURL *models.BatchURL) (*models.BatchURL, error) {
	return nil, nil
}
