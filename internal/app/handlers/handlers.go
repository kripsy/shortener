package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"net/http"

	"github.com/kripsy/shortener/internal/app/models"
	"github.com/kripsy/shortener/internal/app/utils"
	"go.uber.org/zap"
)

type Repository interface {
	CreateOrGetFromStorage(ctx context.Context, url string) (string, error)
	GetOriginalURLFromStorage(ctx context.Context, url string) (string, error)
	CreateOrGetBatchFromStorage(ctx context.Context, batchURL *models.BatchURL) (*models.BatchURL, error)

	Close()
	Ping() error
}

type APIHandler struct {
	repository Repository
	globalURL  string
	myLogger   *zap.Logger
}

func APIHandlerInit(repository Repository, globalURL string, myLogger *zap.Logger) (*APIHandler, error) {

	ht := &APIHandler{
		repository: repository,
		globalURL:  globalURL,
		myLogger:   myLogger,
	}
	return ht, nil
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	val, err := h.repository.CreateOrGetFromStorage(ctx, string(body))
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	url, err := h.repository.GetOriginalURLFromStorage(ctx, shortURL)

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

	h.myLogger.Debug("start SaveURLJSONHandler")
	if r.Method != http.MethodPost || r.Header.Get("Content-Type") != "application/json" {
		h.myLogger.Debug("Bad req", zap.String("Content-Type", r.Header.Get("Content-Type")),
			zap.String("Method", r.Method))
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	var payload URLRequestType
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.myLogger.Debug("Empty body")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &payload)

	if err != nil {
		h.myLogger.Debug("Error unmarshall body", zap.String("error unmarshall", err.Error()))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	h.myLogger.Debug("Unmarshall body", zap.Any("body", payload))

	val, err := h.repository.CreateOrGetFromStorage(context.Background(), payload.URL)
	if err != nil {
		h.myLogger.Debug("Error CreateOrGetFromStorage", zap.String("error CreateOrGetFromStorage", err.Error()))
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(URLResponseType{
		Result: utils.ReturnURL(val, h.globalURL),
	})

	if err != nil {
		h.myLogger.Debug("Error Marshall response", zap.String("error Marshall response", err.Error()))
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

/*
	SaveBatchURLHandler — save batch original url

[

	{
	    "correlation_id": "<строковый идентификатор>",
	    "original_url": "<URL для сокращения>"
	},
	...

]

return
[

	{
	    "correlation_id": "<строковый идентификатор из объекта запроса>",
	    "short_url": "<результирующий сокращённый URL>"
	},
	...

]
*/
func (h *APIHandler) SaveBatchURLHandler(w http.ResponseWriter, r *http.Request) {

	h.myLogger.Debug("start SaveBatchURLHandler")
	if r.Method != http.MethodPost || r.Header.Get("Content-Type") != "application/json" {
		h.myLogger.Debug("Bad req", zap.String("Content-Type", r.Header.Get("Content-Type")),
			zap.String("Method", r.Method))
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	h.myLogger.Debug("Read body", zap.Any("msg", string(body)))
	if err != nil {
		h.myLogger.Debug("Empty body")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var payload *models.BatchURL
	err = json.Unmarshal(body, &payload)
	fmt.Println(len(*payload))
	if err != nil {
		h.myLogger.Debug("Error unmarshall body", zap.String("error unmarshall", err.Error()))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if len(*payload) < 1 {
		h.myLogger.Debug("Payload size < 1")
		http.Error(w, "Empty payload", http.StatusBadRequest)
		return
	}

	h.myLogger.Debug("Unmarshall body", zap.Any("body", payload))

	val, err := h.repository.CreateOrGetBatchFromStorage(context.Background(), payload)
	if err != nil {
		h.myLogger.Debug("Error CreateOrGetFromStorage", zap.String("error CreateOrGetFromStorage", err.Error()))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	// important!!! short url must include server address. It's easy, but in 12 increment i forgot about it
	for k := range *val {
		(*val)[k].ShortURL = utils.ReturnURL((*val)[k].ShortURL, h.globalURL)
	}

	h.myLogger.Debug("Result CreateOrGetBatchFromStorage", zap.Any("msg", val))

	resp, err := json.Marshal(val)

	if err != nil {
		h.myLogger.Debug("Error Marshall response", zap.String("error Marshall response", err.Error()))
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

// PingDBHandler — handler to check success db connection
func (h *APIHandler) PingDBHandler(w http.ResponseWriter, r *http.Request) {

	h.myLogger.Debug("start PingDBHandler")
	err := h.repository.Ping()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}
