// Package logger provides functionality for logger.
package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInitLog(t *testing.T) {
	type args struct {
		level string
	}
	tests := []struct {
		name    string
		args    args
		want    *zap.Logger
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "incorrect log level",
			args: args{
				level: "good level",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "correct log level",
			args: args{
				level: "WARN",
			},
			want:    &zap.Logger{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := InitLog(tt.args.level)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, logger)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, logger)
			}
		})
	}
}
