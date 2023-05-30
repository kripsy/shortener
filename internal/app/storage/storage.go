package storage

import (
	"fmt"

	"github.com/kripsy/shortener/internal/app/logger"
	"github.com/kripsy/shortener/internal/app/utils"
)

type Storage struct {
	storage map[string]string
}

func InitStorageFromFile(storage map[string]string) error {
	events, err := readURL()
	if err != nil {
		logger.Log.Warn("error read URLs")
		return err
	}
	for _, v := range events {
		storage[v.OriginalUrl] = v.ShortUrl
	}
	return nil
}

func InitStorage(initValue map[string]string) *Storage {
	m := Storage{}
	m.storage = initValue
	InitStorageFromFile(m.storage)

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

		e := make([]Event, 1)

		e[0] = *NewEvent(val, url)
		addURL(e)

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
