package logger

import (
	"fmt"

	"go.uber.org/zap"
)

// var Log *zap.Logger = zap.NewNop()

func InitLog(level string) (*zap.Logger, error) {

	lvl, err := zap.ParseAtomicLevel(level)

	if err != nil {
		fmt.Errorf("Failed to parse zap logger level: %s", err.Error())
		return nil, err
	}

	cfg := zap.NewProductionConfig()

	cfg.Level = lvl
	zl, err := cfg.Build()

	if err != nil {
		fmt.Errorf("Failed to build zap logger: %s", err.Error())
		return nil, err
	}

	log := zl
	log.Info(`Logger level`, zap.String("logLevel", level))
	fmt.Println(log)
	return log, nil
}
