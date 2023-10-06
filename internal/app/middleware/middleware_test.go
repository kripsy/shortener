package middleware_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kripsy/shortener/internal/app/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGrpcRequestLogger(t *testing.T) {
	core, recorded := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	m := &middleware.MyMiddleware{
		MyLogger: logger,
	}

	mockHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, fmt.Errorf("%w", status.Error(codes.Internal, "test error"))
	}

	_, err := m.GrpcRequestLogger(context.Background(),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test/method"},
		mockHandler)
	assert.Error(t, err)

	entries := recorded.All()
	assert.Len(t, entries, 1)
	assert.Contains(t, entries[0].Message, "got incoming gRPC request")
}

func TestCompressMiddleware(t *testing.T) {
	logger, err := zap.NewProduction()
	assert.NoError(t, err)

	m := &middleware.MyMiddleware{
		MyLogger: logger,
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("test response"))
		if err != nil {
			return
		}
	})

	compressedHandler := m.CompressMiddleware(handler)

	t.Run("no compression", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "http://example.com", strings.NewReader("test request data"))
		rr := httptest.NewRecorder()

		compressedHandler.ServeHTTP(rr, req)

		assert.Equal(t, "test response", rr.Body.String())
	})

	t.Run("accept gzip with POST request", func(t *testing.T) {
		var buf bytes.Buffer
		writer := gzip.NewWriter(&buf)
		_, err := writer.Write([]byte("test request data"))
		assert.NoError(t, err)
		writer.Close()
		// Создаем POST-запрос с данными в теле
		req := httptest.NewRequest(http.MethodPost, "http://example.com", &buf)
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Content-Encoding", "gzip")
		rr := httptest.NewRecorder()

		compressedHandler.ServeHTTP(rr, req)

		// Проверяем, что ответ содержит заголовок "Content-Encoding: gzip"
		assert.Equal(t, "gzip", rr.Header().Get("Content-Encoding"))

		// Проверяем, что тело ответа сжато
		reader, err := gzip.NewReader(rr.Body)
		assert.NoError(t, err)

		decompressed, err := io.ReadAll(reader)
		assert.NoError(t, err)

		assert.Equal(t, "test response", string(decompressed))
	})
}
