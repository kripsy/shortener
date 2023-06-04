package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kripsy/shortener/internal/app/logger"
	"github.com/kripsy/shortener/internal/app/server"
	"github.com/kripsy/shortener/internal/app/storage"

	"github.com/kripsy/shortener/internal/app/config"
)

func main() {

	config := config.InitConfig()
	myLogger, err := logger.InitLog(config.LoggerLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer myLogger.Sync() // flushes buffer, if any

	fs := storage.InitFileStorageFile(config.FileStoragePath)
	repo := storage.InitStorage(map[string]string{}, fs, myLogger)
	s := server.InitServer(config.URLPrefixRepo, repo, myLogger)
	fmt.Printf("SERVER_ADDRESS: %s\n", config.URLServer)
	fmt.Printf("BASE_URL: %s\n", config.URLPrefixRepo)
	err = http.ListenAndServe(config.URLServer, s.Router)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
