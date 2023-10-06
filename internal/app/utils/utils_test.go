package utils_test

import (
	"context"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/kripsy/shortener/internal/app/utils"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func BenchmarkCreateShortURL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = utils.CreateShortURL()
		// fmt.Println(res)
	}
}

func BenchmarkCreateShortURLWithoutFmt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = utils.CreateShortURLWithoutFmt()
		// fmt.Println(res)
	}
}

func TestGetTokenFromBearer(t *testing.T) {
	type args struct {
		bearerString string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "valid bearer string",
			args: args{
				bearerString: "Bearer validToken",
			},
			want:    "validToken",
			wantErr: false,
		},
		{
			name: "missing bearer string",
			args: args{
				bearerString: "oooppppssss",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "missing token after bearer string",
			args: args{
				bearerString: "Bearer ",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "empty bearer string",
			args: args{
				bearerString: "",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := utils.GetTokenFromBearer(tt.args.bearerString)
			if tt.wantErr {
				assert.Empty(t, token)
				assert.NotEmpty(t, err)
			} else {
				assert.Equal(t, token, tt.want)
				assert.Empty(t, err)
			}
		})
	}
}

func TestGetToken(t *testing.T) {
	type args struct {
		r *http.Request
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "token in Authorization header",
			args: args{
				r: &http.Request{
					Header: http.Header{
						"Authorization": []string{"Bearer validToken"},
					},
				},
			},
			want:    "validToken",
			wantErr: false,
		},
		{
			name: "token in cookie",
			args: args{
				r: &http.Request{
					Header: http.Header{},
				},
			},
			want:    "validToken",
			wantErr: false,
		},
		{
			name: "token not found",
			args: args{
				r: &http.Request{
					Header: http.Header{},
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	tests[1].args.r.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "validToken",
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := utils.GetToken(tt.args.r)
			if tt.wantErr {
				assert.NotEmpty(t, err)
			} else {
				assert.Empty(t, err)
			}
			assert.Equal(t, token, tt.want)
		})
	}
}

func TestStingContains(t *testing.T) {
	type args struct {
		arrayString  []string
		searchString string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "substring in string",
			args: args{
				arrayString:  []string{"test", "test2", "secret"},
				searchString: "secret",
			},
			want: true,
		},
		{
			name: "substring not in string",
			args: args{
				arrayString:  []string{"test", "test2", "secret"},
				searchString: "top secret",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := utils.StingContains(tt.args.arrayString, tt.args.searchString)
			assert.Equal(t, res, tt.want)
		})
	}
}

func TestGetTokenFromMetadata(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		expectedToken := "test_token"
		md := metadata.Pairs("authorization", "Bearer "+expectedToken)
		ctx := metadata.NewIncomingContext(context.Background(), md)

		token, err := utils.GetTokenFromMetadata(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedToken, token)
	})

	t.Run("no metadata", func(t *testing.T) {
		ctx := context.Background()

		_, err := utils.GetTokenFromMetadata(ctx)
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "not metadata in context"))
	})

	t.Run("no token", func(t *testing.T) {
		md := metadata.Pairs("other_header", "value")
		ctx := metadata.NewIncomingContext(context.Background(), md)

		_, err := utils.GetTokenFromMetadata(ctx)
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "token not found"))
	})
}

func TestCreateCertificate(t *testing.T) {
	if _, err := os.Stat("./cert"); os.IsNotExist(err) {
		err := os.Mkdir("./cert", 0755)
		assert.NoError(t, err)
	}
	defer os.Remove(utils.ServerCertPath)
	defer os.Remove(utils.PrivateKeyPath)
	defer os.RemoveAll("./cert")
	err := utils.CreateCertificate()
	assert.NoError(t, err)

	serverCertInfo, err := os.Stat(utils.ServerCertPath)
	assert.NoError(t, err)
	assert.NotZero(t, serverCertInfo.Size())

	privateKeyInfo, err := os.Stat(utils.PrivateKeyPath)
	assert.NoError(t, err)
	assert.NotZero(t, privateKeyInfo.Size())
}

func TestReturnURL(t *testing.T) {
	tests := []struct {
		endpoint  string
		globalURL string
		expected  string
	}{
		{
			endpoint:  "testEndpoint",
			globalURL: "http://example.com",
			expected:  "http://example.com/testEndpoint",
		},
		{
			endpoint:  "anotherEndpoint",
			globalURL: "https://example.org",
			expected:  "https://example.org/anotherEndpoint",
		},
		{
			endpoint:  "",
			globalURL: "https://example.org",
			expected:  "https://example.org/",
		},
		{
			endpoint:  "testEndpoint",
			globalURL: "",
			expected:  "/testEndpoint",
		},
	}

	for _, tt := range tests {
		result := utils.ReturnURL(tt.endpoint, tt.globalURL)
		assert.Equal(t, tt.expected, result)
	}
}

func TestCreateShortURLWithoutFmt(t *testing.T) {
	shortURL, err := utils.CreateShortURLWithoutFmt()
	assert.NoError(t, err)
	assert.Equal(t, 10, len(shortURL))
	anotherShortURL, err := utils.CreateShortURLWithoutFmt()
	assert.NoError(t, err)
	assert.NotEqual(t, shortURL, anotherShortURL)
}

func TestCreateShortURL(t *testing.T) {
	shortURL, err := utils.CreateShortURL()
	assert.NoError(t, err)
	assert.Equal(t, 10, len(shortURL))
	anotherShortURL, err := utils.CreateShortURL()
	assert.NoError(t, err)
	assert.NotEqual(t, shortURL, anotherShortURL)
}
