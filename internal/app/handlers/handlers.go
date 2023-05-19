package handlers

import (
	"io"

	"net/http"

	"github.com/kripsy/shortener/internal/app/utils"
)

type Repository interface {
	CreateOrGetFromStorage(url string) (string, error)
	GetFromStorage(url string) (string, error)
}

type APIHandler struct {
	storage   Repository
	globalURL string
}

func APIHandlerInit(storage Repository, globalURL string) *APIHandler {
	ht := &APIHandler{
		storage:   storage,
		globalURL: globalURL,
	}
	return ht
}

// SaveURLHandler — save original url, create short url into storage
func (h *APIHandler) SaveURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	val, err := h.storage.CreateOrGetFromStorage(string(body))
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, utils.ReturnURL(val, h.globalURL))
}

// GetURLHandler — get origin url from storage by shortURL
func (h *APIHandler) GetURLHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	// remove first slash
	shortURL := (r.URL.Path)[1:]

	url, err := h.storage.GetFromStorage(shortURL)

	// if we got error in getFromStorage - bad request
	if err != nil {

		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "plain/text")
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)

}
