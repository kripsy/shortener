package main

import (
	"fmt"
	"os"

	"github.com/kripsy/shortener/internal/app/config"
	"github.com/kripsy/shortener/internal/app/logger"
	"github.com/kripsy/shortener/internal/app/server"
	"github.com/kripsy/shortener/internal/app/storage"
	"go.uber.org/zap"
)

type app struct {
	Config      *config.Config
	MyLogger    *zap.Logger
	FileStorage *storage.FileStorage
	Storage     *storage.Storage
	Server      *server.MyServer
}

func newApp() (*app, error) {

	config := config.InitConfig()
	myLogger, err := logger.InitLog(config.LoggerLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return nil, err
	}

	fs := storage.InitFileStorageFile(config.FileStoragePath)
	s := storage.InitStorage(map[string]string{}, fs, myLogger)
	srv := server.InitServer(config.URLPrefixRepo, s, myLogger)
	return &app{
		Config:      config,
		MyLogger:    myLogger,
		FileStorage: fs,
		Storage:     s,
		Server:      srv,
	}, nil
}
