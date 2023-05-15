package mymemory

import (
	"fmt"

	"github.com/kripsy/shortener/internal/app/utils"
)

type MyMemory struct {
	myMemory map[string]string
}

func InitMyMemory(initValue map[string]string) *MyMemory {
	m := MyMemory{}
	m.myMemory = initValue

	return &m
}

func (m *MyMemory) CreateOrGetFromMemory(url string) (string, error) {
	val, ok := m.myMemory[string(url)]
	// If the key exists
	if !ok {
		// input into our memory
		val, err := utils.CreateShortURL()
		if err != nil {
			return "", err
		}
		m.myMemory[string(url)] = val
		return val, nil
	}
	return val, nil
}

func (m MyMemory) GetFromMemory(url string) (string, error) {

	var val string
	ok := false
	// for every key from MYMEMORY check our shortURL. If exist set `val = k` and `ok = true`

	for k, v := range m.myMemory {
		if v == string(url) {

			ok = true
			val = k
		}
	}
	// If the key exists
	if ok {
		return val, nil
	}
	// key not exist
	return "", fmt.Errorf("not exists")
}
