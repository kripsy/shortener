package server

import (
	"github.com/gorilla/mux"
	"github.com/kripsy/shortener/internal/app/handlers"
)

type MyServer struct {
	Router    *mux.Router
	ServerUrl string
}

func InitServer(serverUrl string, repo handlers.Repository) *MyServer {
	m := MyServer{}
	m.ServerUrl = serverUrl
	m.Router = mux.NewRouter()
	m.Router.HandleFunc(`/`, handlers.SaveUrlHandler(repo, m.ServerUrl))
	m.Router.HandleFunc(`/{id}`, handlers.GetUrlHandler(repo, m.ServerUrl))

	return &m
}
