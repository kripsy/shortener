package utils

import (
	"crypto/rand"
	"fmt"
)

func CreateShortURL() (string, error) {
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

func ReturnURL(endpoint, globalURL string) string {
	return globalURL + "/" + endpoint
}

func StingContains(arrayString []string, searchString string) bool {
	for _, v := range arrayString {
		if v == searchString {
			return true
		}
	}
	return false
}
