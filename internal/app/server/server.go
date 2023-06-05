package server

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/kripsy/shortener/internal/app/handlers"
	"github.com/kripsy/shortener/internal/app/middleware"
)

type MyServer struct {
	Router   *chi.Mux
	URLRepo  string
	MyLogger *zap.Logger
}

func InitServer(URLRepo string, repo handlers.Repository, myLogger *zap.Logger) *MyServer {
	m := &MyServer{}
	m.URLRepo = URLRepo
	m.MyLogger = myLogger
	ht := handlers.APIHandlerInit(repo, m.URLRepo, m.MyLogger)
	myMiddleware := middleware.InitMyMiddleware(m.MyLogger)
	m.Router = chi.NewRouter()
	m.Router.Use(myMiddleware.CompressMiddleware)
	m.Router.Use(myMiddleware.RequestLogger)
	m.Router.Post("/", ht.SaveURLHandler)
	m.Router.Get("/{id}", ht.GetURLHandler)
	m.Router.Post("/api/shorten", ht.SaveURLJSONHandler)

	return m
}
