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
		start := time.Now()
		// Код для измерения
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		logger.Log.Debug("got incoming HTTP request",
			zap.String("URI", r.URL.String()),
			zap.String("method", r.Method),
			zap.Int64("duration (Nanoseconds)", duration.Nanoseconds()),
		)
	})
}
