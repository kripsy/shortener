package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	//nolint:depguard
	"github.com/kripsy/shortener/internal/app/handlers"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type MyMiddleware struct {
	MyLogger      *zap.Logger
	repo          handlers.Repository
	TrustedSubnet *net.IPNet
}

var (
	//nolint:gochecknoglobals
	once sync.Once
	//nolint:gochecknoglobals
	instance *MyMiddleware
)

func InitMyMiddleware(myLogger *zap.Logger, repo handlers.Repository, ts *net.IPNet) *MyMiddleware {
	once.Do(func() {
		instance = &MyMiddleware{
			MyLogger:      myLogger,
			repo:          repo,
			TrustedSubnet: ts,
		}
	})

	return instance
}

// RequestLogger — middleware-логер для входящих HTTP-запросов.
func (m *MyMiddleware) RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// start duration work of handler
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w,            // insert orig http.ResponseWriter
			responseData:   responseData, // insert our respData
		}

		next.ServeHTTP(&lw, r)
		duration := time.Since(start)
		m.MyLogger.Info("got incoming HTTP request",
			zap.String("URI", r.URL.String()),
			zap.String("method", r.Method),
			zap.Int64("duration (Nanoseconds)", duration.Nanoseconds()),
			zap.Int64("status code", int64(responseData.status)),
			zap.Int64("response size", int64(responseData.size)),
		)
	})
}

// CompressMiddleware — middleware for compress and decompress data.
func (m *MyMiddleware) CompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w
		m.MyLogger.Debug("start CompressMiddleware")
		if r.Header.Get("Content-Encoding") != "gzip" {
			m.MyLogger.Debug("continue without compress")
			next.ServeHTTP(ow, r)

			return
		}
		m.MyLogger.Debug("continue with compress")
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			m.MyLogger.Debug("Accept-Encoding gzip")
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			m.MyLogger.Debug("Content-Encoding gzip")
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				return
			}
			// меняем тело запроса на новое
			r.Body = cr
			defer cr.Close()
		}
		next.ServeHTTP(ow, r)
	})
}

// RequestLogger — middleware-логер для входящих HTTP-запросов.
func (m *MyMiddleware) GrpcRequestLogger(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	// Запуск отсчета времени выполнения
	start := time.Now()

	// Вызов следующего обработчика в цепочке
	resp, err := handler(ctx, req)

	// Вычисление продолжительности выполнения
	duration := time.Since(start)

	// Получение статуса ответа
	status, _ := status.FromError(err)

	// Логирование информации о запросе
	m.MyLogger.Info("got incoming gRPC request",
		zap.String("method", info.FullMethod),
		zap.Int64("duration (Nanoseconds)", duration.Nanoseconds()),
		zap.String("status code", status.Code().String()),
		// Дополнительные поля можно добавить по необходимости
	)

	return resp, err
}
