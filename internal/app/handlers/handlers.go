// Package handlers provides handlers for web server.
package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"net/http"

	"github.com/kripsy/shortener/internal/app/auth"
	"github.com/kripsy/shortener/internal/app/models"
	"github.com/kripsy/shortener/internal/app/usecase"
	"github.com/kripsy/shortener/internal/app/utils"
	"go.uber.org/zap"
)

type Repository interface {
	CreateOrGetFromStorage(ctx context.Context, url string, userID int) (string, error)
	GetOriginalURLFromStorage(ctx context.Context, url string) (string, error)
	CreateOrGetBatchFromStorage(ctx context.Context, batchURL *models.BatchURL, userID int) (*models.BatchURL, error)
	RegisterUser(ctx context.Context) (*models.User, error)
	GetBatchURLFromStorage(ctx context.Context, userID int) (*models.BatchURL, error)
	DeleteSliceURLFromStorage(ctx context.Context, shortURL []string, userID int) error
	GetStatsFromStorage(ctx context.Context) (*models.Stats, error)

	GetUserByID(ctx context.Context, id int) (*models.User, error)
	Close()
	Ping() error
}

type APIHandler struct {
	myLogger   *zap.Logger
	repository Repository
	globalURL  string
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

// SaveURLHandler — save original url, create short url into storage.
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
	isUniqueError := false

	token, _ := utils.GetToken(r)
	userID, _ := auth.GetUserID(token)
	//nolint:contextcheck
	val, err := h.repository.CreateOrGetFromStorage(ctx, string(body), userID)
	if err != nil {
		var ue *models.UniqueError
		if errors.As(err, &ue) {
			isUniqueError = true
		} else {
			h.myLogger.Debug("Error CreateOrGetFromStorage", zap.String("error CreateOrGetFromStorage", err.Error()))
			http.Error(w, "", http.StatusBadRequest)

			return
		}
	}
	w.Header().Set("Content-Type", "plain/text")
	if isUniqueError {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	_, err = io.WriteString(w, utils.ReturnURL(val, h.globalURL))
	if err != nil {
		h.myLogger.Debug("Error CreateOrGetFromStorage", zap.String("error WriteString", err.Error()))
		http.Error(w, "", http.StatusInternalServerError)

		return
	}
}

// GetURLHandler — get origin url from storage by shortURL.
func (h *APIHandler) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusBadRequest)

		return
	}
	// remove first slash
	shortURL := (r.URL.Path)[1:]
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	//nolint:contextcheck
	url, err := h.repository.GetOriginalURLFromStorage(ctx, shortURL)
	fmt.Println(url)
	// if we got error in getFromStorage - bad request
	if err != nil {
		var isDeletedError *models.IsDeletedError
		if errors.As(err, &isDeletedError) {
			h.myLogger.Debug("URL is deleted", zap.String("msg", shortURL))
			http.Error(w, "", http.StatusGone)

			return
		}
		http.Error(w, "", http.StatusBadRequest)

		return
	}

	w.Header().Set("Content-Type", "plain/text")
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// SaveURLJSONHandler — save original url, create short url into storage with JSON.
func (h *APIHandler) SaveURLJSONHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := utils.GetToken(r)
	userID, _ := auth.GetUserID(token)
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
	isUniqueError := false
	//nolint:contextcheck
	val, err := h.repository.CreateOrGetFromStorage(context.Background(), payload.URL, userID)
	if err != nil {
		var ue *models.UniqueError
		if errors.As(err, &ue) {
			isUniqueError = true
		} else {
			h.myLogger.Debug("Error CreateOrGetFromStorage", zap.String("error CreateOrGetFromStorage", err.Error()))
			http.Error(w, "", http.StatusBadRequest)

			return
		}
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
	if isUniqueError {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	_, err = w.Write(resp)
	if err != nil {
		h.myLogger.Debug("Error CreateOrGetFromStorage", zap.String("error Write", err.Error()))
		http.Error(w, "", http.StatusInternalServerError)

		return
	}
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

].
*/
func (h *APIHandler) SaveBatchURLHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := utils.GetToken(r)
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

	batch := &models.BatchURL{}
	err = json.Unmarshal(body, batch)
	if err != nil {
		h.myLogger.Debug("error unmarshall to body", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)

		return
	}

	result, err := usecase.ProcessBatchURLs(r.Context(), batch, h.repository, token, h.globalURL, h.myLogger)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
	resp, err := json.Marshal(result)

	if err != nil {
		h.myLogger.Debug("Error Marshall response", zap.String("error Marshall response", err.Error()))
		http.Error(w, "", http.StatusBadRequest)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(resp)
	if err != nil {
		h.myLogger.Debug("Error CreateOrGetBatchFromStorage", zap.String("error Write", err.Error()))
		http.Error(w, "", http.StatusInternalServerError)

		return
	}
}

/*
GetBatchURLHandler - handler, that return all urls, that user have sent.
If Batch is empty - return 204.
*/
func (h *APIHandler) GetBatchURLHandler(w http.ResponseWriter, r *http.Request) {
	h.myLogger.Debug("start GetBatchURLHandler")
	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusBadRequest)

		return
	}
	token, _ := utils.GetToken(r)
	fmt.Println(token)
	userID, _ := auth.GetUserID(token)
	fmt.Println(userID)
	//nolint:contextcheck
	batchURL, err := h.repository.GetBatchURLFromStorage(context.Background(), userID)
	// if we got error in getFromStorage - bad request
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	// if len batchURL 0 - send 204
	if len(*batchURL) == 0 {
		w.WriteHeader(http.StatusNoContent)

		return
	}

	for i := range *batchURL {
		(*batchURL)[i].ShortURL = utils.ReturnURL((*batchURL)[i].ShortURL, h.globalURL)
	}

	resp, err := json.Marshal(batchURL)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)

		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		h.myLogger.Debug("Error GetBatchURLHandler", zap.String("error Write", err.Error()))
		http.Error(w, "", http.StatusInternalServerError)

		return
	}
}

// PingDBHandler — handler to check success db connection.
func (h *APIHandler) PingDBHandler(w http.ResponseWriter, _ *http.Request) {
	h.myLogger.Debug("start PingDBHandler")
	err := h.repository.Ping()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)

		return
	}
	w.WriteHeader(http.StatusOK)
}
func (h *APIHandler) DeleteBatchURLHandler(w http.ResponseWriter, r *http.Request) {
	h.myLogger.Debug("start GetBatchURLHandler")
	if r.Method != http.MethodDelete {
		http.Error(w, "", http.StatusBadRequest)

		return
	}
	token, _ := utils.GetToken(r)
	fmt.Println(token)
	userID, _ := auth.GetUserID(token)
	fmt.Println(userID)

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		h.myLogger.Debug("Empty body")
		http.Error(w, "", http.StatusInternalServerError)

		return
	}
	h.myLogger.Debug("Read body", zap.Any("msg", string(body)))

	str := string(body)

	str = strings.ReplaceAll(str, `[`, "")
	str = strings.ReplaceAll(str, `]`, "")
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, `"`, "")
	splitstr := strings.Split(str, ",")
	h.myLogger.Debug("Result body", zap.Any("msg", splitstr))
	if splitstr[0] == "" {
		h.myLogger.Debug("Bad req, splitstr[0] is empty")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		h.myLogger.Debug("goroutine started with urls", zap.Any("msg", splitstr))
		err = h.repository.DeleteSliceURLFromStorage(ctx, splitstr, userID)
		if err != nil {
			h.myLogger.Debug("error in goroutine DeleteSliceURLFromStorage", zap.String("msg", err.Error()))

			return
		}
	}()

	w.WriteHeader(http.StatusAccepted)
}

/*
GetInternalStats - handler, that returns stats about urls and users.
If X-Real-IP not in trusted_subnet - it will return 403.
*/
func (h *APIHandler) GetInternalStats(w http.ResponseWriter, r *http.Request) {
	h.myLogger.Debug("start GetInternalStats")
	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusBadRequest)

		return
	}

	// realIP := r.Header.Get("X-Real-IP")
	// h.myLogger.Debug("X-Real-IP: ", zap.String("msg", realIP))
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	stats, err := h.repository.GetStatsFromStorage(ctx)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}

	resp, err := json.Marshal(stats)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}
