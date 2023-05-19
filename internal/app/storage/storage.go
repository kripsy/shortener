package storage

import (
	"fmt"

	"github.com/kripsy/shortener/internal/app/utils"
)

type Storage struct {
	storage map[string]string
}

func InitStorage(initValue map[string]string) *Storage {
	m := Storage{}
	m.storage = initValue

	return &m
}

func (m *Storage) CreateOrGetFromStorage(url string) (string, error) {
	// If the key exists
	val, ok := m.storage[string(url)]
	if !ok {
		// input into our storage
		val, err := utils.CreateShortURL()
		if err != nil {
			return "", err
		}
		m.storage[string(url)] = val
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
