package utils

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// func for generation random short URL, consist of 5 bytes
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

// concatination global URL (address server) and endpoint (short URL)
func ReturnURL(endpoint, globalURL string) string {
	return globalURL + "/" + endpoint
}

// function for check is searchString contain in arrayString
func StingContains(arrayString []string, searchString string) bool {
	for _, v := range arrayString {
		if v == searchString {
			return true
		}
	}
	return false
}

// function to get token from header that consist "Bearer ...."
func GetTokenFromBearer(bearerString string) (string, error) {
	splitString := strings.Split(bearerString, "Bearer ")
	fmt.Printf("len splitString %d\n", len(splitString))
	if len(splitString) < 2 {
		fmt.Printf("bearer string not valid")
		return "", fmt.Errorf("bearer string not valid")
	}
	tokenString := splitString[1]
	if tokenString == "" {
		fmt.Printf("tokenString is empty")
		return "", fmt.Errorf("tokenString is empty")
	}
	return tokenString, nil
}

func GetToken(w http.ResponseWriter, r *http.Request) (string, error) {
	var token string
	// try get token from header
	tokenString := r.Header.Get("Authorization")
	if tokenString != "" {
		fmt.Printf("get token from header: %s\n", tokenString)
		token, _ = GetTokenFromBearer(tokenString)
		fmt.Printf("token %s\n", token)
	}
	if token != "" {
		return token, nil
	}

	// if we continue - it means that in header isn't token. Try find it in cookie
	cookieToken, err := r.Cookie("token")
	if err != nil {
		return "", errors.Wrap(err, "cannot get token from cookie")
	}
	token = cookieToken.Value
	if token != "" {
		return token, nil
	}
	return "", fmt.Errorf("token not found")

}
