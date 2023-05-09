package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_saveUrlHandler(t *testing.T) {

	mymemory := map[string]string{
		"https://google.com/": "http://localhost:8080/82643f4619",
	}

	type want struct {
		contentType string
		statusCode  int
	}

	tests := []struct {
		name       string
		request    string
		body       string
		methodType string
		mymemory   map[string]string
		want       want
	}{
		// TODO: Add test cases.
		{
			name:       "First success save originalUrl",
			request:    "/",
			methodType: http.MethodPost,
			body:       "https://practicum.yandex.ru/",
			mymemory:   mymemory,
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
			mymemory:   mymemory,
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
			h := http.HandlerFunc(saveUrlHandler(mymemory))
			h(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			if result.StatusCode == 201 {
				assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
				shortUrl, err := ioutil.ReadAll(result.Body)
				require.NoError(t, err)

				err = result.Body.Close()
				require.NoError(t, err)

				assert.NotEmpty(t, shortUrl, "shortUrl is empty, but except not empty")
			}
		})
	}
}

func Test_getUrlHandler(t *testing.T) {
	mymemory := map[string]string{
		"https://google.com/": "82643f4619",
	}

	type want struct {
		statusCode int
		Location   string
	}

	tests := []struct {
		name       string
		request    string
		body       string
		methodType string
		mymemory   map[string]string
		want       want
	}{
		// TODO: Add test cases.
		{
			name:       "Success get originalUrl",
			request:    "/82643f4619",
			methodType: http.MethodGet,
			mymemory:   mymemory,
			want: want{
				statusCode: 307,
				Location:   "https://google.com/",
			},
		},
		{
			name:       "No success get originalUrl",
			request:    "/82643f4610",
			methodType: http.MethodGet,
			mymemory:   mymemory,
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
			h := http.HandlerFunc(getUrlHandler(mymemory))
			h(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			if result.StatusCode == http.StatusTemporaryRedirect {
				assert.Equal(t, tt.want.Location, result.Header.Get("Location"))
			}
		})
	}
}
