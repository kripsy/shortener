package main

import (
	"fmt"
	"net/http"

	"github.com/kripsy/shortener/internal/app/config"
	"github.com/kripsy/shortener/internal/app/mymemory"
	"github.com/kripsy/shortener/internal/app/server"
)

func main() {

	config := config.InitConfig()
	repo := mymemory.InitMyMemory(map[string]string{})
	s := server.InitServer(config.URLPrefixRepo, repo)
	fmt.Printf("SERVER_ADDRESS: %s\n", config.URLServer)
	fmt.Printf("BASE_URL: %s\n", config.URLPrefixRepo)
	err := http.ListenAndServe(config.URLServer, s.Router)
	if err != nil {
		panic(err)
	}
}
