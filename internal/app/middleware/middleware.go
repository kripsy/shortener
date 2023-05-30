package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/kripsy/shortener/internal/app/logger"
	"go.uber.org/zap"
)

// RequestLogger — middleware-логер для входящих HTTP-запросов.
func RequestLogger(next http.Handler) http.Handler {
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
		logger.Log.Info("got incoming HTTP request",
			zap.String("URI", r.URL.String()),
			zap.String("method", r.Method),
			zap.Int64("duration (Nanoseconds)", duration.Nanoseconds()),
			zap.Int64("status code", int64(responseData.status)),
			zap.Int64("response size", int64(responseData.size)),
		)
	})
}

// CompressMiddleware — middleware for compress and decompress data.
func CompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w
		logger.Log.Debug("start CompressMiddleware")
		if (r.Header.Get("Content-Type") != "application/json") && (r.Header.Get("Content-Type") == "text/html") {
			logger.Log.Debug("continue without compress")
			next.ServeHTTP(ow, r)
			return
		}
		logger.Log.Debug("continue with compress")
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			logger.Log.Debug("Accept-Encoding gzip")
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			logger.Log.Debug("Content-Encoding gzip")
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
