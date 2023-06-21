package inmemorystorage

import (
	"context"
	"fmt"

	"github.com/kripsy/shortener/internal/app/auth"
	"github.com/kripsy/shortener/internal/app/models"
	"github.com/kripsy/shortener/internal/app/utils"
	"go.uber.org/zap"
)

type InMemoryStorage struct {
	storage  map[string]string
	myLogger *zap.Logger
}

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
	m.myLogger.Debug("Start CreateOrGetBatchFromStorage")
	for k, v := range *batchURL {
		shortURL, err := m.CreateOrGetFromStorage(context.Background(), v.OriginalURL)
		if err != nil {
			return nil, err
		}
		(*batchURL)[k].ShortURL = shortURL
		(*batchURL)[k].OriginalURL = ""
	}
	return batchURL, nil
}

func (m InMemoryStorage) GetUserByID(ctx context.Context, ID uint64) (*auth.User, error) {
	return nil, fmt.Errorf("not implemented")
}
