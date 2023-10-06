// middleware_test for unit test protect handler by ip
package middleware_test

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kripsy/shortener/internal/app/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestTrustedSubnetMiddleware(t *testing.T) {
	logger, _ := zap.NewProduction()
	_, ts, _ := net.ParseCIDR("192.168.1.0/24")
	m := &middleware.MyMiddleware{
		MyLogger: logger,
		// Здесь вы можете указать ваш доверенный подсет
		TrustedSubnet: ts,
	}

	handler := m.TrustedSubnetMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	tests := []struct {
		name           string
		url            string
		ip             string
		expectedStatus int
	}{
		{
			name:           "Unprotected URL",
			url:            "/api/some-other-endpoint",
			ip:             "192.168.1.1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Protected URL with trusted IP",
			url:            "/api/internal/stats",
			ip:             "192.168.1.1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Protected URL with untrusted IP",
			url:            "/api/internal/stats",
			ip:             "10.0.0.1",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, tt.url, nil)
			req.Header.Set("X-Real-IP", tt.ip)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
