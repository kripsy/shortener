// middleware for protect handler by ip
package middleware

import (
	"net"
	"net/http"

	//nolint:depguard

	"github.com/kripsy/shortener/internal/app/utils"
	"go.uber.org/zap"
)

// TrustedSubnetMiddleware implements protection the URL depending on the specified IP address.
func (m *MyMiddleware) TrustedSubnetMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		protectedURL := []string{
			"/api/internal/stats",
		}
		m.MyLogger.Debug("Start TrustedSubnetMiddleware")

		// check if current URL is protected
		isURLProtected := utils.StingContains(protectedURL, r.URL.Path)
		m.MyLogger.Debug("URL protected value", zap.Bool("msg", isURLProtected))

		if !isURLProtected {
			next.ServeHTTP(w, r)

			return
		}

		// try get X-Real-IP from request header
		ip := net.ParseIP(r.Header.Get("X-Real-IP"))

		if (m.TrustedSubnet != nil) && m.TrustedSubnet.Contains(ip) {
			m.MyLogger.Debug("empty X-Real-IP")
			next.ServeHTTP(w, r)
		} else {
			m.MyLogger.Debug("X-Real-IP not in trusted subnet")
			w.WriteHeader(http.StatusForbidden)

			return
		}
	})
}
