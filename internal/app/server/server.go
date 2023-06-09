package server

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/kripsy/shortener/internal/app/handlers"
	"github.com/kripsy/shortener/internal/app/middleware"
)

type MyServer struct {
	Router   *chi.Mux
	MyDB     handlers.Repository
	MyLogger *zap.Logger
	URLRepo  string
}

func InitServer(URLRepo string, repo handlers.Repository, myLogger *zap.Logger) (*MyServer, error) {
	m := &MyServer{
		Router:   chi.NewRouter(),
		MyLogger: myLogger,
		URLRepo:  URLRepo,
	}

	ht, err := handlers.APIHandlerInit(repo, m.URLRepo, m.MyLogger)
	if err != nil {
		return nil, err
	}
	myMiddleware := middleware.InitMyMiddleware(m.MyLogger)

	m.Router.Use(myMiddleware.CompressMiddleware)
	m.Router.Use(myMiddleware.RequestLogger)
	m.Router.Post("/", ht.SaveURLHandler)
	m.Router.Get("/{id}", ht.GetURLHandler)
	m.Router.Post("/api/shorten", ht.SaveURLJSONHandler)
	m.Router.Get("/ping", ht.PingDBHandler)

	return m, nil
}
