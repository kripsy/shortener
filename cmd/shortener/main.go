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
	s := server.InitServer(config.UrlPrefixRepo, repo)
	fmt.Printf("SERVER_ADDRESS: %s\n", config.UrlServer)
	fmt.Printf("BASE_URL: %s\n", config.UrlPrefixRepo)
	err := http.ListenAndServe(config.UrlServer, s.Router)
	if err != nil {
		panic(err)
	}
}
