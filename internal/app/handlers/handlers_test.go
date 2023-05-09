package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kripsy/shortener/internal/app/mymemory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveURLHandler(t *testing.T) {

	myMemory := mymemory.InitMyMemory(map[string]string{})

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

		myMemory Repository

		want want
	}{
		// TODO: Add test cases.
		{
			name:       "First success save originalUrl",
			request:    "/",
			methodType: http.MethodPost,
			body:       "https://practicum.yandex.ru/",
			myMemory:   myMemory,
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
			myMemory:   myMemory,
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

			h := http.HandlerFunc(SaveURLHandler(tt.myMemory, globalURL))

			h(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			if result.StatusCode == 201 {
				assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
				shortURL, err := ioutil.ReadAll(result.Body)
				require.NoError(t, err)

				err = result.Body.Close()
				require.NoError(t, err)

				assert.NotEmpty(t, shortURL, "shortURL is empty, but except not empty")
			}
		})
	}
}

func TestGetURLHandler(t *testing.T) {
	myMemory := mymemory.InitMyMemory(map[string]string{
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
		myMemory   Repository

		want want
	}{
		// TODO: Add test cases.
		{
			name:       "Success get originalUrl",
			request:    "/82643f4619",
			methodType: http.MethodGet,
			myMemory:   myMemory,
			want: want{
				statusCode: 307,
				Location:   "https://google.com/",
			},
		},
		{
			name:       "No success get originalUrl",
			request:    "/82643f4610",
			methodType: http.MethodGet,
			myMemory:   myMemory,
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
			h := http.HandlerFunc(GetURLHandler(tt.myMemory, globalURL))

			h(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			if result.StatusCode == http.StatusTemporaryRedirect {
				assert.Equal(t, tt.want.Location, result.Header.Get("Location"))
			}
		})
	}
}
