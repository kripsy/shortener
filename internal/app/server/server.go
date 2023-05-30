package server

import (
	"github.com/go-chi/chi/v5"

	"github.com/kripsy/shortener/internal/app/handlers"
	"github.com/kripsy/shortener/internal/app/middleware"
)

type MyServer struct {
	Router  *chi.Mux
	URLRepo string
}

func InitServer(URLRepo string, repo handlers.Repository) *MyServer {
	m := &MyServer{}
	m.URLRepo = URLRepo
	ht := handlers.APIHandlerInit(repo, m.URLRepo)
	m.Router = chi.NewRouter()
	m.Router.Use(middleware.CompressMiddleware)
	m.Router.Use(middleware.RequestLogger)
	m.Router.Post("/", ht.SaveURLHandler)
	m.Router.Get("/{id}", ht.GetURLHandler)
	m.Router.Post("/api/shorten", ht.SaveURLJSONHandler)

	return m
}
