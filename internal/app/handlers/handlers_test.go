package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"

	"github.com/kripsy/shortener/internal/app/logger"
	"github.com/kripsy/shortener/internal/app/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// t.Setenv("SERVER_ADDRESS", "localhost:8080")
// t.Setenv("BASE_URL", "http://localhost:8080")
// t.Setenv("LOG_LEVEL", "Info")
// t.Setenv("DATABASE_DSN", "")
// t.Setenv("FILE_STORAGE_PATH", "")
// envPrefixAddr := os.Getenv("BASE_URL")
type TestParams struct {
	TestLogger     *zap.Logger
	TestPrefixAddr string
}

func getParamsForTest() *TestParams {
	tl, _ := logger.InitLog("Debug")

	tp := &TestParams{
		TestLogger:     tl,
		TestPrefixAddr: "http://localhost:8080",
	}
	return tp
}

func TestSaveURLHandler(t *testing.T) {

	paramTest := getParamsForTest()

	type want struct {
		contentType string
		statusCode  int
	}

	tests := []struct {
		name       string
		body       string
		methodType string
		want       want
	}{
		// TODO: Add test cases.
		{
			name:       "First success save originalUrl",
			methodType: http.MethodPost,
			body:       "https://practicum.yandex.ru/",
			want: want{
				contentType: "plain/text",
				statusCode:  201,
			},
		},
		{
			name:       "No success save originalUrl",
			methodType: http.MethodGet,
			body:       "https://practicum.yandex.ru/123",
			want: want{
				statusCode: 400,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mdb := mocks.NewMockRepository(ctrl)
			mdb.EXPECT().CreateOrGetFromStorage(gomock.Any(), gomock.Any()).Return("good", nil).AnyTimes()

			body := strings.NewReader(tt.body)

			request := httptest.NewRequest(tt.methodType, "/", body)
			w := httptest.NewRecorder()
			ht, _ := APIHandlerInit(mdb, paramTest.TestPrefixAddr, paramTest.TestLogger)
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
	paramTest := getParamsForTest()
	type want struct {
		statusCode int
		Location   string
	}

	tests := []struct {
		name       string
		request    string
		body       string
		methodType string
		want       want
	}{
		// TODO: Add test cases.
		{
			name:       "Success get originalUrl",
			request:    "/82643f4619",
			methodType: http.MethodGet,
			want: want{
				statusCode: 307,
				Location:   "https://google.com/",
			},
		},
		{
			name:       "No success get originalUrl",
			request:    "/82643f4610",
			methodType: http.MethodGet,
			want: want{
				statusCode: 400,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mdb := mocks.NewMockRepository(ctrl)
			mdb.EXPECT().GetOriginalURLFromStorage(gomock.Any(), "82643f4619").Return("https://google.com/", nil).AnyTimes()
			mdb.EXPECT().GetOriginalURLFromStorage(gomock.Any(), "82643f4610").Return("", errors.New("1")).AnyTimes()

			body := strings.NewReader(tt.body)
			request := httptest.NewRequest(tt.methodType, tt.request, body)
			w := httptest.NewRecorder()
			ht, _ := APIHandlerInit(mdb, paramTest.TestPrefixAddr, paramTest.TestLogger)
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

func TestSaveURLJSONHandler(t *testing.T) {
	paramTest := getParamsForTest()

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

		want want
	}{
		{
			// TODO: Add test cases.
			name:        "First success save originalUrl",
			request:     "/",
			methodType:  http.MethodPost,
			body:        `{"url":"123"}`,
			contentType: "application/json",

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

			want: want{
				contentType: "application/json",
				statusCode:  400,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mdb := mocks.NewMockRepository(ctrl)
			mdb.EXPECT().CreateOrGetFromStorage(gomock.Any(), gomock.Any()).Return("good", nil).AnyTimes()
			body := strings.NewReader(tt.body)

			request := httptest.NewRequest(tt.methodType, tt.request, body)
			request.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			ht, _ := APIHandlerInit(mdb, paramTest.TestPrefixAddr, paramTest.TestLogger)
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

func TestAPIHandler_PingDBHandler(t *testing.T) {
	paramTest := getParamsForTest()

	type want struct {
		statusCode int
	}

	tests := []struct {
		name        string
		request     string
		body        string
		methodType  string
		contentType string
		success     bool
		want        want
	}{
		// TODO: Add test cases.
		{
			// TODO: Add test cases.
			name:        "First success ping",
			request:     "/ping",
			methodType:  http.MethodGet,
			success:     true,
			contentType: "application/json",
			want: want{
				statusCode: 200,
			},
		},
		{
			// TODO: Add test cases.
			name:        "First failed ping",
			request:     "/ping",
			methodType:  http.MethodGet,
			success:     false,
			contentType: "application/json",
			want: want{
				statusCode: 500,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mdb := mocks.NewMockRepository(ctrl)

			if tt.success {
				mdb.EXPECT().Ping().Return(nil).AnyTimes()
			} else {
				mdb.EXPECT().Ping().Return(errors.New("test")).AnyTimes()
			}

			request := httptest.NewRequest(tt.methodType, tt.request, nil)
			request.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()

			ht, _ := APIHandlerInit(mdb, paramTest.TestPrefixAddr, paramTest.TestLogger)
			h := ht.PingDBHandler

			h(w, request)
			result := w.Result()
			require.Equal(t, tt.want.statusCode, result.StatusCode)
			result.Body.Close()
		})
	}
}
