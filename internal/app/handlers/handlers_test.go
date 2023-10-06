package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kripsy/shortener/internal/app/handlers"
	"github.com/kripsy/shortener/internal/app/logger"
	"github.com/kripsy/shortener/internal/app/mocks"
	"github.com/kripsy/shortener/internal/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type TestParams struct {
	testLogger     *zap.Logger
	testPrefixAddr string
}

func getParamsForTest() *TestParams {
	tl, _ := logger.InitLog("Debug")

	tp := &TestParams{
		testLogger:     tl,
		testPrefixAddr: "http://localhost:8080",
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
			mdb.EXPECT().CreateOrGetFromStorage(gomock.Any(), gomock.Any(), gomock.Any()).Return("good", nil).AnyTimes()

			body := strings.NewReader(tt.body)

			request := httptest.NewRequest(tt.methodType, "/", body)
			w := httptest.NewRecorder()
			ht, _ := handlers.APIHandlerInit(mdb, paramTest.testPrefixAddr, paramTest.testLogger)
			h := ht.SaveURLHandler

			h(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			if result.StatusCode == http.StatusCreated {
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
		Location   string
		statusCode int
	}

	tests := []struct {
		name       string
		request    string
		body       string
		methodType string
		want       want
	}{
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
			//nolint:goerr113
			mdb.EXPECT().GetOriginalURLFromStorage(gomock.Any(), "82643f4610").Return("", errors.New("1")).AnyTimes()

			body := strings.NewReader(tt.body)
			request := httptest.NewRequest(tt.methodType, tt.request, body)
			w := httptest.NewRecorder()
			ht, _ := handlers.APIHandlerInit(mdb, paramTest.testPrefixAddr, paramTest.testLogger)
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
			//
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
			//
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
			mdb.EXPECT().CreateOrGetFromStorage(gomock.Any(), gomock.Any(), gomock.Any()).Return("good", nil).AnyTimes()
			body := strings.NewReader(tt.body)

			request := httptest.NewRequest(tt.methodType, tt.request, body)
			request.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			ht, _ := handlers.APIHandlerInit(mdb, paramTest.testPrefixAddr, paramTest.testLogger)
			h := ht.SaveURLJSONHandler

			h(w, request)
			result := w.Result()

			shortURL, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			err = result.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			if result.StatusCode != http.StatusCreated {
				return
			}

			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			var resp handlers.URLResponseType
			err = json.Unmarshal(shortURL, &resp)
			require.NoError(t, err)

			assert.NotEmpty(t, resp.Result, "Orig URL is empty, but except not empty")
		})
	}
}

func TestPingDBHandler(t *testing.T) {
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
		{
			//
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
			//
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
				//nolint:goerr113
				mdb.EXPECT().Ping().Return(errors.New("test")).AnyTimes()
			}

			request := httptest.NewRequest(tt.methodType, tt.request, nil)
			request.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()

			ht, _ := handlers.APIHandlerInit(mdb, paramTest.testPrefixAddr, paramTest.testLogger)
			h := ht.PingDBHandler

			h(w, request)
			result := w.Result()
			require.Equal(t, tt.want.statusCode, result.StatusCode)
			result.Body.Close()
		})
	}
}

func TestSaveBatchURLHandler(t *testing.T) {
	paramTest := getParamsForTest()

	type want struct {
		contentType string
		body        string
		statusCode  int
	}

	tests := []struct {
		name        string
		request     string
		body        string
		methodType  string
		contentType string
		want        want
	}{
		{
			//
			name:       "First success save originalUrl",
			request:    "/",
			methodType: http.MethodPost,
			body: `[
				{
					"correlation_id": "1",
					"original_url":   "https://ya.ru"
				},
				{
					"correlation_id": "2",
					"original_url":   "https://google.com"
				}
			]`,
			contentType: "application/json",

			want: want{
				contentType: "application/json",
				statusCode:  201,
				body: `[
					{
						"correlation_id": "1",
						"short_url":   "ttttt"
					},
					{
						"correlation_id": "2",
						"short_url":   "fffff"
					}
				]`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mdb := mocks.NewMockRepository(ctrl)
			var valueInut *models.BatchURL
			var valueOutput *models.BatchURL
			err := json.Unmarshal([]byte(tt.body), &valueInut)
			assert.NoError(t, err)
			err = json.Unmarshal([]byte(tt.want.body), &valueOutput)
			assert.NoError(t, err)
			fmt.Println(valueInut)
			fmt.Println(valueOutput)

			mdb.EXPECT().CreateOrGetBatchFromStorage(gomock.Any(), valueInut, gomock.Any()).Return(valueOutput, nil).AnyTimes()
			//nolint:noctx
			request, err := http.NewRequest(tt.methodType, tt.request, strings.NewReader(tt.body))
			assert.NoError(t, err)

			request.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			ht, _ := handlers.APIHandlerInit(mdb, paramTest.testPrefixAddr, paramTest.testLogger)
			h := ht.SaveBatchURLHandler

			h(w, request)

			result := w.Result()
			defer result.Body.Close()

			resp, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			var respModel *models.BatchURL
			err = json.Unmarshal(resp, &respModel)
			assert.NoError(t, err)
			assert.Equal(t, respModel, valueOutput)
		})
	}
}

func TestDeleteBatchURLHandler(t *testing.T) {
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
		want        want
	}{
		{
			//
			name:        "First success delete urls",
			request:     "/",
			methodType:  http.MethodDelete,
			body:        `["9260e84518", "bf04361ccc", "a75fee90c6", "f04361ccc"]`,
			contentType: "application/json",

			want: want{
				contentType: "application/json",
				statusCode:  202,
			},
		},
		{
			name:        "First fail to delete urls",
			request:     "/",
			methodType:  http.MethodDelete,
			body:        ``,
			contentType: "application/json",

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
			mdb.EXPECT().DeleteSliceURLFromStorage(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			//nolint:noctx
			request, err := http.NewRequest(tt.methodType, tt.request, strings.NewReader(tt.body))
			assert.NoError(t, err)

			request.Header.Set("Content-Type", tt.contentType)

			w := httptest.NewRecorder()
			ht, _ := handlers.APIHandlerInit(mdb, paramTest.testPrefixAddr, paramTest.testLogger)
			h := ht.DeleteBatchURLHandler
			h(w, request)

			result := w.Result()
			defer result.Body.Close()

			_, err = io.ReadAll(request.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestGetInternalStats(t *testing.T) {
	paramTest := getParamsForTest()

	type want struct {
		contentType string
		statusCode  int
	}

	tests := []struct {
		name        string
		request     string
		methodType  string
		contentType string
		want        want
	}{
		{
			name:        "First success get stats",
			request:     "/",
			methodType:  http.MethodGet,
			contentType: "application/json",
			want: want{
				contentType: "application/json",
				statusCode:  200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mdb := mocks.NewMockRepository(ctrl)
			mdb.EXPECT().GetStatsFromStorage(gomock.Any()).Return(&models.Stats{
				URLs:  10,
				Users: 10,
			}, nil)

			request, err := http.NewRequestWithContext(context.Background(), tt.methodType, tt.request, nil)
			assert.NoError(t, err, "error create request")

			w := httptest.NewRecorder()
			ht, err := handlers.APIHandlerInit(mdb, paramTest.testPrefixAddr, paramTest.testLogger)
			assert.NoError(t, err, "error create api handler")

			h := ht.GetInternalStats
			h(w, request)

			result := w.Result()

			defer result.Body.Close()

			resp, err := io.ReadAll(result.Body)
			assert.NoError(t, err, "error read body response")

			if tt.want.statusCode == http.StatusOK {
				var stats models.Stats
				err = json.Unmarshal(resp, &stats)
				assert.NoError(t, err, "error unmarshall resp")
			}
		})
	}
}
