// middleware for protect handler by ip
package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"

	//nolint:depguard

	"github.com/kripsy/shortener/internal/app/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	grpcstatus "google.golang.org/grpc/status"
)

// TrustedSubnetMiddleware implements protection the URL depending on the specified IP address.
func (m *MyMiddleware) TrustedSubnetMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.MyLogger.Debug("Start TrustedSubnetMiddleware")
		isURLProtected := m.urlIsProtected(r.URL.Path)
		m.MyLogger.Debug("URL protected value", zap.Bool("msg", isURLProtected))
		if !isURLProtected {
			next.ServeHTTP(w, r)

			return
		}

		// try get X-Real-IP from request header
		ip := net.ParseIP(r.Header.Get("X-Real-IP"))

		m.MyLogger.Debug("X-Real-IP", zap.Any("msg", ip))
		if m.isIPTrusted(ip) {
			m.MyLogger.Debug("empty X-Real-IP")
			next.ServeHTTP(w, r)
		} else {
			m.MyLogger.Debug("X-Real-IP not in trusted subnet")
			w.WriteHeader(http.StatusForbidden)

			return
		}
	})
}

// TrustedSubnetMiddleware implements protection the URL depending on the specified IP address.
func (m *MyMiddleware) GrpcTrustedSubnetMiddleware(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	m.MyLogger.Debug("Start TrustedSubnetMiddleware")

	methodName := info.FullMethod
	m.MyLogger.Debug("method ", zap.String("msg", methodName))
	if methodName != "/Shortener/Stats" {
		return handler(ctx, req)
	}
	peerInfo, ok := peer.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("%w", grpcstatus.Error(codes.PermissionDenied, "No peer info found"))
	}
	ip, _, err := net.SplitHostPort(peerInfo.Addr.String())
	if err != nil {
		return nil, fmt.Errorf("%w", grpcstatus.Error(codes.Internal, "Failed to parse IP address"))
	}

	if m.isIPTrusted(net.ParseIP(ip)) {
		return handler(ctx, req)
	}

	return nil, fmt.Errorf("%w", grpcstatus.Error(codes.PermissionDenied, "IP not trusted"))
}

func (m *MyMiddleware) isIPTrusted(ip net.IP) bool {
	if m.TrustedSubnet != nil && m.TrustedSubnet.Contains(ip) {
		return true
	}

	return false
}

func (m *MyMiddleware) urlIsProtected(url string) bool {
	protectedURL := []string{
		"/api/internal/stats",
	}
	m.MyLogger.Debug("Start GrpcTrustedSubnetMiddleware")
	urlIsProtected := utils.StingContains(protectedURL, url)
	m.MyLogger.Debug("URL protected value", zap.Bool("msg", urlIsProtected))
	// check if current URL is protected
	return urlIsProtected
}
