package handlers

import (
	"io"

	"net/http"

	"github.com/kripsy/shortener/internal/app/utils"
)

type Repository interface {
	CreateOrGetFromMemory(url string) (string, error)
	GetFromMemory(url string) (string, error)
}

type HandlerType struct {
	myMemory  Repository
	globalURL string
}

func HandlerTypeInit(myMemory Repository, globalURL string) *HandlerType {
	ht := &HandlerType{
		myMemory:  myMemory,
		globalURL: globalURL,
	}
	return ht
}

// SaveURLHandler — save original url, create short url into memory
func (h *HandlerType) SaveURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	val, err := h.myMemory.CreateOrGetFromMemory(string(body))
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, utils.ReturnURL(val, h.globalURL))
}

// GetURLHandler — get origin url from memory by shortURL
func (h *HandlerType) GetURLHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	// remove first slash
	shortURL := (r.URL.Path)[1:]

	url, err := h.myMemory.GetFromMemory(shortURL)

	// if we got error in getFromMemory - bad request
	if err != nil {

		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "plain/text")
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)

}
