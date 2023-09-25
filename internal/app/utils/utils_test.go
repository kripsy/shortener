package utils_test

import (
	"net/http"
	"testing"

	"github.com/kripsy/shortener/internal/app/utils"
	"github.com/stretchr/testify/assert"
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
