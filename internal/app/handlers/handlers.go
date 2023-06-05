package handlers

import (
	"encoding/json"
	"io"

	"net/http"

	"github.com/kripsy/shortener/internal/app/utils"
	"go.uber.org/zap"
)

type Repository interface {
	CreateOrGetFromStorage(url string) (string, error)
	GetFromStorage(url string) (string, error)
}

type APIHandler struct {
	storage   Repository
	globalURL string
	MyLogger  *zap.Logger
}

func APIHandlerInit(storage Repository, globalURL string, myLogger *zap.Logger) *APIHandler {
	ht := &APIHandler{
		storage:   storage,
		globalURL: globalURL,
		MyLogger:  myLogger,
	}
	return ht
}

type URLResponseType struct {
	Result string `json:"result,omitempty"`
}

type URLRequestType struct {
	URL string `json:"url,omitempty"`
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

// SaveURLHandler — save original url, create short url into storage with JSON
func (h *APIHandler) SaveURLJSONHandler(w http.ResponseWriter, r *http.Request) {

	h.MyLogger.Debug("start SaveURLJSONHandler")
	if r.Method != http.MethodPost || r.Header.Get("Content-Type") != "application/json" {
		h.MyLogger.Debug("Bad req", zap.String("Content-Type", r.Header.Get("Content-Type")),
			zap.String("Method", r.Method))
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	var payload URLRequestType
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.MyLogger.Debug("Empty body")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &payload)

	if err != nil {
		h.MyLogger.Debug("Error unmarshall body", zap.String("error unmarshall", err.Error()))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	h.MyLogger.Debug("Unmarshall body", zap.Any("body", payload))

	val, err := h.storage.CreateOrGetFromStorage(payload.URL)
	if err != nil {
		h.MyLogger.Debug("Error CreateOrGetFromStorage", zap.String("error CreateOrGetFromStorage", err.Error()))
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(URLResponseType{
		Result: utils.ReturnURL(val, h.globalURL),
	})

	if err != nil {
		h.MyLogger.Debug("Error Marshall response", zap.String("error Marshall response", err.Error()))
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}
