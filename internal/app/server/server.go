package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/kripsy/shortener/internal/app/handlers"
)

type MyServer struct {
	Router  *chi.Mux
	UrlRepo string
}

func InitServer(UrlRepo string, repo handlers.Repository) *MyServer {
	m := MyServer{}
	m.UrlRepo = UrlRepo

	m.Router = chi.NewRouter()
	m.Router.Post("/", handlers.SaveUrlHandler(repo, m.UrlRepo))
	m.Router.Get("/{id}", handlers.GetUrlHandler(repo, m.UrlRepo))

	return &m
}
