package middleware

import (
	"net/http"
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
