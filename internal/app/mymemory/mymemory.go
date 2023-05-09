package mymemory

import (
	"crypto/rand"
	"fmt"
)

type MyMemory struct {
	myMemory map[string]string
}

func InitMyMemory(initValue map[string]string) *MyMemory {
	m := MyMemory{}
	m.myMemory = initValue

	return &m
}

func (m *MyMemory) CreateOrGetFromMemory(url []byte) (string, error) {
	val, ok := m.myMemory[string(url)]
	// If the key exists
	if ok {
		return val, nil
	}
	// input into our memory
	if val, err := createShortURL(url); err == nil {
		m.myMemory[string(url)] = val
		return val, nil
	} else {
		return "", err
	}
}

func (m MyMemory) GetFromMemory(url []byte) (string, error) {

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

func createShortURL(input []byte) (string, error) {
	// create slice 5 bytes
	buf := make([]byte, 5)

	// call rand.Read.
	_, err := rand.Read(buf)

	// if error - return empty string and error
	if err != nil {
		return "", fmt.Errorf("error while generating random string: %s", err)
	}

	// print bytes in hex and return as string
	return fmt.Sprintf("%x", buf), nil
}
