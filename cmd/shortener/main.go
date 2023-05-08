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

func safeUrl(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	val, err := createOrGetFromMemory(body)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, returnUrl(val))
}

func getUrl(w http.ResponseWriter, r *http.Request) {
	// remove first slash
	shortUrl := (r.URL.Path)[1:]
	url, err := getFromMemory([]byte(shortUrl))

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

func main() {

	r := mux.NewRouter()
	r.HandleFunc(`/`, safeUrl).Methods("POST")
	r.HandleFunc(`/{id}`, getUrl).Methods("GET")

	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		panic(err)
	}

}

func returnUrl(endpoint string) string {
	return URL + "/" + endpoint
}

func createOrGetFromMemory(url []byte) (string, error) {
	val, ok := MYMEMORY[string(url)]
	// If the key exists
	if ok {
		return val, nil
	}
	// input into our memory
	if val, err := createShortUrl(url); err == nil {
		MYMEMORY[string(url)] = val
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

func getFromMemory(url []byte) (string, error) {
	var val string
	ok := false
	// for every key from MYMEMORY check our shortUrl. If exist set `val = k` and `ok = true`
	for k, v := range MYMEMORY {
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
