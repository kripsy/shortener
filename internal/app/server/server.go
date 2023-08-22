package server

import (
	"expvar"
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/kripsy/shortener/internal/app/handlers"
	"github.com/kripsy/shortener/internal/app/middleware"
)

type MyServer struct {
	Router   *chi.Mux
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
	myMiddleware := middleware.InitMyMiddleware(m.MyLogger, repo)

	m.Router.Use(myMiddleware.CompressMiddleware)
	m.Router.Use(myMiddleware.RequestLogger)
	m.Router.Use(myMiddleware.JWTMiddleware)
	m.Router.Post("/", ht.SaveURLHandler)
	m.Router.Get("/{id}", ht.GetURLHandler)
	m.Router.Post("/api/shorten", ht.SaveURLJSONHandler)
	m.Router.Post("/api/shorten/batch", ht.SaveBatchURLHandler)
	m.Router.Get("/ping", ht.PingDBHandler)
	m.Router.Get("/api/user/urls", ht.GetBatchURLHandler)
	m.Router.Delete("/api/user/urls", ht.DeleteBatchURLHandler)

	m.Router.HandleFunc("/debug/pprof/*", pprof.Index)
	m.Router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	m.Router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	m.Router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	m.Router.HandleFunc("/debug/pprof/trace", pprof.Trace)
	m.Router.HandleFunc("/debug/vars", expVars)

	m.Router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	m.Router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	m.Router.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	m.Router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	m.Router.Handle("/debug/pprof/block", pprof.Handler("block"))
	m.Router.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))

	return m, nil
}

// Replicated from expvar.go as not public.
func expVars(w http.ResponseWriter, r *http.Request) {
	first := true
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\n")
	expvar.Do(func(kv expvar.KeyValue) {
		if !first {
			fmt.Fprintf(w, ",\n")
		}
		first = false
		fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)
	})
	fmt.Fprintf(w, "\n}\n")
}
