package application

import (
	"context"
	"fmt"
	"os"

	"github.com/kripsy/shortener/internal/app/config"
	database "github.com/kripsy/shortener/internal/app/db"
	"github.com/kripsy/shortener/internal/app/filestorage"
	"github.com/kripsy/shortener/internal/app/handlers"
	"github.com/kripsy/shortener/internal/app/inmemorystorage"
	"github.com/kripsy/shortener/internal/app/logger"
	"github.com/kripsy/shortener/internal/app/models"
	"github.com/kripsy/shortener/internal/app/server"
	"go.uber.org/zap"
)

type App struct {
	appConfig *config.Config
	appLogger *zap.Logger
	appServer *server.MyServer
	appRepo   handlers.Repository
}

func (a *App) GetAppServer() *server.MyServer {
	return a.appServer
}

func (a *App) GetAppConfig() *config.Config {
	return a.appConfig
}

func (a *App) GetAppRepo() handlers.Repository {
	return a.appRepo
}

func (a *App) GetAppLogger() *zap.Logger {
	return a.appLogger
}

func NewApp(ctx context.Context) (*App, error) {
	cfg := config.InitConfig()
	var repo handlers.Repository

	myLogger, err := logger.InitLog(cfg.LoggerLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return nil, err
	}

	switch cfg.RepositoryType {
	case config.PostgresDB:
		var db *database.PostgresDB
		db, err = database.InitDB(cfg.DatabaseDsn, myLogger)
		if err != nil {
			myLogger.Debug("Failed init DB", zap.String("msg", err.Error()))
			return nil, err
		}

		repo = db

	case config.FileStorage:
		var fs *filestorage.FileStorage
		fs, err = filestorage.InitFileStorageFile(cfg.FileStoragePath, myLogger)
		if err != nil {
			myLogger.Debug("Failed init filestorage", zap.String("msg", err.Error()))
			return nil, err
		}
		repo = fs

	case config.InMemory:
		var inmemory *inmemorystorage.InMemoryStorage
		inmemory, err = inmemorystorage.InitInMemoryStorage(map[string]models.Event{}, myLogger)
		if err != nil {
			myLogger.Debug("Failed init inmemorystorage", zap.String("msg", err.Error()))
			return nil, err
		}
		repo = inmemory
	}

	srv, err := server.InitServer(cfg.URLPrefixRepo, repo, myLogger)

	if err != nil {
		myLogger.Debug("Failed init server", zap.String("msg", err.Error()))
		return nil, err
	}

	return &App{
		appConfig: cfg,
		appLogger: myLogger,
		appServer: srv,
		appRepo:   repo,
	}, nil
}
