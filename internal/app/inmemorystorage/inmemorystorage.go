// Package inmemorystorage provides functionality inmemory storage.
package inmemorystorage

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/kripsy/shortener/internal/app/models"
	"github.com/kripsy/shortener/internal/app/utils"
	"go.uber.org/zap"
)

type InMemoryStorage struct {
	storage  map[string]models.Event
	myLogger *zap.Logger
	rwmutex  *sync.RWMutex
}

func InitInMemoryStorage(initValue map[string]models.Event, myLogger *zap.Logger) (*InMemoryStorage, error) {
	m := &InMemoryStorage{
		storage:  initValue,
		myLogger: myLogger,
		rwmutex:  &sync.RWMutex{},
	}

	return m, nil
}

func (m *InMemoryStorage) CreateOrGetFromStorageWithoutPointer(_ context.Context,
	url string,
	userID int) (string, error) {
	// If the key exists
	m.rwmutex.RLock()
	val, ok := m.storage[url]
	m.rwmutex.RUnlock()

	if !ok {
		// input into our storage
		val, err := utils.CreateShortURL()
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}

		event := models.NewEventWithoutPointer(val, url, userID)
		m.rwmutex.Lock()
		defer m.rwmutex.Unlock()
		m.storage[url] = event

		return val, nil
	}

	return val.ShortURL, nil
}

func (m *InMemoryStorage) CreateOrGetFromStorage(_ context.Context, url string, userID int) (string, error) {
	// If the key exists
	m.rwmutex.RLock()
	val, ok := m.storage[url]
	m.rwmutex.RUnlock()
	if !ok {
		// input into our storage
		val, err := utils.CreateShortURL()
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}

		event := models.NewEvent(val, url, userID)
		m.rwmutex.Lock()
		defer m.rwmutex.Unlock()
		m.storage[url] = *event

		return val, nil
	}

	return val.ShortURL, nil
}

func (m InMemoryStorage) GetOriginalURLFromStorage(_ context.Context, shortURL string) (string, error) {
	event := &models.Event{}
	ok := false
	// for every key from MYMEMORY check our shortURL. If exist set `val = k` and `ok = true`
	m.rwmutex.RLock()
	defer m.rwmutex.RUnlock()
	for k, v := range m.storage {
		if v.ShortURL == shortURL {
			ok = true
			event.OriginalURL = k

			break
		}
	}
	if !ok {
		// key not exist
		//nolint:goerr113,nolintlint
		return "", fmt.Errorf("not exists")
	}
	// If the key exists
	return event.OriginalURL, nil
}

func (m InMemoryStorage) Close() {
}

func (m InMemoryStorage) Ping() error {
	return nil
}

func (m InMemoryStorage) CreateOrGetBatchFromStorage(ctx context.Context,
	batchURL *models.BatchURL,
	_ int) (*models.BatchURL, error) {
	m.myLogger.Debug("Start CreateOrGetBatchFromStorage")
	for k, v := range *batchURL {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		defer cancel()
		shortURL, err := m.CreateOrGetFromStorage(ctx, v.OriginalURL, 1)
		if err != nil {
			return nil, err
		}
		(*batchURL)[k].ShortURL = shortURL
		(*batchURL)[k].OriginalURL = ""
	}

	return batchURL, nil
}

func (m InMemoryStorage) GetUserByID(_ context.Context, id int) (*models.User, error) {
	m.rwmutex.RLock()
	defer m.rwmutex.RUnlock()
	for _, v := range m.storage {
		if v.UserID == id {
			return &models.User{ID: id}, nil
		}
	}
	//nolint:goerr113
	return nil, fmt.Errorf("user not found")
}

func (m InMemoryStorage) RegisterUser(_ context.Context) (*models.User, error) {
	return &models.User{
		ID: int(uuid.New().ID()),
	}, nil
}

func (m InMemoryStorage) GetBatchURLFromStorage(_ context.Context, userID int) (*models.BatchURL, error) {
	batchURL := &models.BatchURL{}
	m.rwmutex.RLock()
	defer m.rwmutex.RUnlock()
	for _, v := range m.storage {
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

func (m InMemoryStorage) DeleteSliceURLFromStorage(_ context.Context, _ []string, _ int) error {
	fmt.Println("Not implemented yet")

	return nil
}
