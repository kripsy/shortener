package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kripsy/shortener/internal/app/logger"
	"github.com/kripsy/shortener/internal/app/storage"

	"github.com/kripsy/shortener/internal/app/config"
	"github.com/kripsy/shortener/internal/app/server"
)

func main() {
	defer logger.Log.Sync() // flushes buffer, if any

	config := config.InitConfig()
	logger.InitLog(config.LoggerLevel)
	storage.InitFileStorageFile(config.FileStoragePath)
	repo := storage.InitStorage(map[string]string{})
	s := server.InitServer(config.URLPrefixRepo, repo)
	fmt.Printf("SERVER_ADDRESS: %s\n", config.URLServer)
	fmt.Printf("BASE_URL: %s\n", config.URLPrefixRepo)
	err := http.ListenAndServe(config.URLServer, s.Router)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
