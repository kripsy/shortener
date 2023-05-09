package main

import (
	"net/http"

	"github.com/kripsy/shortener/internal/app/config"
	"github.com/kripsy/shortener/internal/app/mymemory"
	"github.com/kripsy/shortener/internal/app/server"
)

func main() {

	config := config.InitConfig("http://localhost", "8080")
	repo := mymemory.InitMyMemory(map[string]string{})
	s := server.InitServer(config.Url+`:`+config.Port, repo)

	err := http.ListenAndServe(`:`+config.Port, s.Router)
	if err != nil {
		panic(err)
	}
}
