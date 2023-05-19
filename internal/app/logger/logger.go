package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger = zap.NewNop()

func InitLog(level string) error {

	lvl, err := zap.ParseAtomicLevel(level)

	if err != nil {
		Log.Fatal("can't parse lovel zap logger", zap.String("message", err.Error()))
		return err
	}

	cfg := zap.NewProductionConfig()

	cfg.Level = lvl
	zl, err := cfg.Build()

	if err != nil {
		Log.Fatal("can't build zap logger", zap.String("message", err.Error()))
		return err
	}

	Log = zl

	Log.Info(`Logger level`, zap.String("logLevel", level))

	return nil
}
