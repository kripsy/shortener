package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

const URL = "http://localhost:8080"

var MYMEMORY map[string]string = map[string]string{}

// saveUrlHandler — save original url, create short url into memory
func saveUrlHandler(myMemory map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		val, err := createOrGetFromMemory(body, myMemory)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "plain/text")
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, returnUrl(val))
	}
}

// getUrlHandler — get origin url from memory by shortUrl
func getUrlHandler(myMemory map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		// remove first slash
		shortUrl := (r.URL.Path)[1:]

		url, err := getFromMemory([]byte(shortUrl), myMemory)

		// if we got error in getFromMemory - bad request
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "plain/text")
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc(`/`, saveUrlHandler(MYMEMORY))    //.Methods("POST")
	r.HandleFunc(`/{id}`, getUrlHandler(MYMEMORY)) //.Methods("GET")

	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		panic(err)
	}

}

func returnUrl(endpoint string) string {
	return URL + "/" + endpoint
}

func createOrGetFromMemory(url []byte, myMemory map[string]string) (string, error) {
	val, ok := myMemory[string(url)]
	// If the key exists
	if ok {
		return val, nil
	}
	// input into our memory
	if val, err := createShortUrl(url); err == nil {
		myMemory[string(url)] = val
		return val, nil
	} else {
		return "", err
	}
}

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

func getFromMemory(url []byte, myMemory map[string]string) (string, error) {
	var val string
	ok := false
	// for every key from MYMEMORY check our shortUrl. If exist set `val = k` and `ok = true`
	for k, v := range myMemory {
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
	return "", fmt.Errorf("Not exists")
}
