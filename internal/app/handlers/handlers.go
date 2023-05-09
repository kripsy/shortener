package handlers

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/kripsy/shortener/internal/app/utils"
)

type Repository interface {
	CreateOrGetFromMemory(url []byte) (string, error)
	GetFromMemory(url []byte) (string, error)
}

// SaveUrlHandler — save original url, create short url into memory
func SaveUrlHandler(myMemory Repository, globalUrl string) http.HandlerFunc {
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

		val, err := myMemory.CreateOrGetFromMemory(body)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "plain/text")
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, utils.ReturnUrl(val, globalUrl))
	}
}

// getUrlHandler — get origin url from memory by shortUrl
func GetUrlHandler(myMemory Repository, globalUrl string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		// remove first slash
		shortUrl := (r.URL.Path)[1:]

		url, err := myMemory.GetFromMemory([]byte(shortUrl))

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
