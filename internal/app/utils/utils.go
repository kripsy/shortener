package utils

import (
	"crypto/rand"
	"fmt"
)

func createShortUrl(input []byte) (string, error) {
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

func ReturnUrl(endpoint, globalUrl string) string {
	return globalUrl + "/" + endpoint
}
