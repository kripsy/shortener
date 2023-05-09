package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/kripsy/shortener/internal/app/handlers"
)

type MyServer struct {
	Router  *chi.Mux
	URLRepo string
}

func InitServer(URLRepo string, repo handlers.Repository) *MyServer {
	m := MyServer{}
	m.URLRepo = URLRepo

	m.Router = chi.NewRouter()
	m.Router.Post("/", handlers.SaveURLHandler(repo, m.URLRepo))
	m.Router.Get("/{id}", handlers.GetURLHandler(repo, m.URLRepo))

	return &m
}
