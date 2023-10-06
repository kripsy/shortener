package middleware

import (
	"context"
	"net"
	"testing"

	"github.com/kripsy/shortener/internal/app/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type TestParams struct {
	testLogger *zap.Logger
}

func getParamsForTest() *TestParams {
	tl, _ := logger.InitLog("Debug")

	tp := &TestParams{
		testLogger: tl,
	}

	return tp
}
func TestIsIPTrusted(t *testing.T) {
	tests := []struct {
		name          string
		trustedSubnet string
		ipToCheck     string
		expected      bool
	}{
		{
			name:          "IP is trusted",
			trustedSubnet: "192.168.1.0/24",
			ipToCheck:     "192.168.1.10",
			expected:      true,
		},
		{
			name:          "IP is not trusted",
			trustedSubnet: "192.168.1.0/24",
			ipToCheck:     "192.168.2.10",
			expected:      false,
		},
		{
			name:          "No trusted subnet set",
			trustedSubnet: "",
			ipToCheck:     "192.168.1.10",
			expected:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, subnet, _ := net.ParseCIDR(tt.trustedSubnet)
			m := &MyMiddleware{
				TrustedSubnet: subnet,
			}
			ip := net.ParseIP(tt.ipToCheck)
			result := m.isIPTrusted(ip)

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGrpcTrustedSubnetMiddleware(t *testing.T) {
	tp := getParamsForTest()
	tests := []struct {
		name          string
		methodName    string
		ip            string
		trustedSubnet string
		expectedError error
	}{
		{
			name:          "Trusted IP for GetStats",
			methodName:    "/shortener.v1.ShortenerService/GetStats",
			ip:            "192.168.1.10",
			trustedSubnet: "192.168.1.0/24",
			expectedError: nil,
		},
		{
			name:          "Untrusted IP for GetStats",
			methodName:    "/shortener.v1.ShortenerService/GetStats",
			ip:            "192.168.2.10",
			trustedSubnet: "192.168.1.0/24",
			expectedError: status.Error(codes.PermissionDenied, "IP not trusted"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, subnet, _ := net.ParseCIDR(tt.trustedSubnet)
			m := &MyMiddleware{
				TrustedSubnet: subnet,
				MyLogger:      tp.testLogger,
			}

			ctx := context.Background()
			addr := &net.TCPAddr{IP: net.ParseIP(tt.ip)}
			p := &peer.Peer{Addr: addr}
			ctx = peer.NewContext(ctx, p)

			info := &grpc.UnaryServerInfo{FullMethod: tt.methodName}
			_, err := m.GrpcTrustedSubnetMiddleware(ctx,
				nil,
				info,
				func(ctx context.Context,
					req interface{}) (interface{}, error) {
					//nolint:nilnil
					return nil, nil
				})

			if tt.expectedError != nil {
				assert.NotEmpty(t, err)
			} else {
				assert.Empty(t, err)
			}
		})
	}
}
