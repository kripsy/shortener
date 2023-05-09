package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/kripsy/shortener/internal/app/handlers"
)

type MyServer struct {
	Router    *chi.Mux
	ServerUrl string
}

func InitServer(serverUrl string, repo handlers.Repository) *MyServer {
	m := MyServer{}
	m.ServerUrl = serverUrl

	m.Router = chi.NewRouter()
	m.Router.Post("/", handlers.SaveUrlHandler(repo, m.ServerUrl))
	m.Router.Get("/", handlers.GetUrlHandler(repo, m.ServerUrl))

	return &m
}
