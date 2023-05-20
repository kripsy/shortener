package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kripsy/shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveURLHandler(t *testing.T) {

	storage := storage.InitStorage(map[string]string{})

	globalURL := "http://localhost:8080"

	type want struct {
		contentType string
		statusCode  int
	}

	tests := []struct {
		name       string
		request    string
		body       string
		methodType string

		storage Repository

		want want
	}{
		// TODO: Add test cases.
		{
			name:       "First success save originalUrl",
			request:    "/",
			methodType: http.MethodPost,
			body:       "https://practicum.yandex.ru/",
			storage:    storage,
			want: want{
				contentType: "plain/text",
				statusCode:  201,
			},
		},
		{
			name:       "No success save originalUrl",
			request:    "/",
			methodType: http.MethodGet,
			body:       "https://practicum.yandex.ru/123",
			storage:    storage,
			want: want{
				statusCode: 400,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := strings.NewReader(tt.body)

			request := httptest.NewRequest(tt.methodType, tt.request, body)
			w := httptest.NewRecorder()
			ht := APIHandlerInit(tt.storage, globalURL)
			h := ht.SaveURLHandler

			h(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			if result.StatusCode == 201 {
				assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
				shortURL, err := io.ReadAll(result.Body)
				require.NoError(t, err)

				err = result.Body.Close()
				require.NoError(t, err)

				assert.NotEmpty(t, shortURL, "shortURL is empty, but except not empty")
			}
		})
	}
}

func TestGetURLHandler(t *testing.T) {
	storage := storage.InitStorage(map[string]string{
		"https://google.com/": "82643f4619",
	})

	globalURL := "http://localhost:8080"
	type want struct {
		statusCode int
		Location   string
	}

	tests := []struct {
		name       string
		request    string
		body       string
		methodType string
		storage    Repository

		want want
	}{
		// TODO: Add test cases.
		{
			name:       "Success get originalUrl",
			request:    "/82643f4619",
			methodType: http.MethodGet,
			storage:    storage,
			want: want{
				statusCode: 307,
				Location:   "https://google.com/",
			},
		},
		{
			name:       "No success get originalUrl",
			request:    "/82643f4610",
			methodType: http.MethodGet,
			storage:    storage,
			want: want{
				statusCode: 400,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := strings.NewReader(tt.body)
			request := httptest.NewRequest(tt.methodType, tt.request, body)
			w := httptest.NewRecorder()
			ht := APIHandlerInit(tt.storage, globalURL)
			h := ht.GetURLHandler

			h(w, request)
			result := w.Result()
			err := result.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			if result.StatusCode == http.StatusTemporaryRedirect {
				assert.Equal(t, tt.want.Location, result.Header.Get("Location"))
			}
		})
	}
}

func TestAPIHandler_SaveURLJSONHandler(t *testing.T) {
	storage := storage.InitStorage(map[string]string{})

	globalURL := "http://localhost:8080"

	type want struct {
		contentType string
		statusCode  int
	}

	tests := []struct {
		name        string
		request     string
		body        string
		methodType  string
		contentType string
		storage     Repository

		want want
	}{
		{
			// TODO: Add test cases.
			name:        "First success save originalUrl",
			request:     "/",
			methodType:  http.MethodPost,
			body:        `{"url":"123"}`,
			contentType: "application/json",
			storage:     storage,
			want: want{
				contentType: "application/json",
				statusCode:  201,
			},
		},
		{
			// TODO: Add test cases.
			name:        "Bad content-type",
			request:     "/",
			methodType:  http.MethodPost,
			body:        `{"url":"123"}`,
			contentType: "text/plain",
			storage:     storage,
			want: want{
				contentType: "application/json",
				statusCode:  400,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := strings.NewReader(tt.body)

			request := httptest.NewRequest(tt.methodType, tt.request, body)
			request.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			ht := APIHandlerInit(tt.storage, globalURL)
			h := ht.SaveURLJSONHandler

			h(w, request)
			result := w.Result()

			shortURL, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			err = result.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			if result.StatusCode != 201 {
				return
			}

			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			var resp URLResponseType
			err = json.Unmarshal(shortURL, &resp)
			require.NoError(t, err)

			assert.NotEmpty(t, resp.Result, "Orig URL is empty, but except not empty")

		})
	}
}
