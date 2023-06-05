package logger

import (
	"go.uber.org/zap"
)

func InitLog(level string) (*zap.Logger, error) {

	lvl, err := zap.ParseAtomicLevel(level)

	if err != nil {
		return nil, err
	}

	cfg := zap.NewProductionConfig()

	cfg.Level = lvl
	zl, err := cfg.Build()

	if err != nil {
		return nil, err

	}

	log := zl
	log.Info(`Logger level`, zap.String("logLevel", level))
	return log, nil
}
