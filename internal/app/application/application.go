package application

import (
	"context"
	"fmt"
	"os"

	"github.com/kripsy/shortener/internal/app/config"
	"github.com/kripsy/shortener/internal/app/db"
	"github.com/kripsy/shortener/internal/app/filestorage"
	"github.com/kripsy/shortener/internal/app/handlers"
	"github.com/kripsy/shortener/internal/app/inmemorystorage"
	"github.com/kripsy/shortener/internal/app/logger"
	"github.com/kripsy/shortener/internal/app/server"
	"go.uber.org/zap"
)

type App struct {
	appConfig *config.Config
	appLogger *zap.Logger
	// FileStorage *storage.FileStorage
	// Storage     *storage.Storage
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
		db, err := db.InitDB(cfg.DatabaseDsn)
		if err != nil && cfg.DatabaseDsn != "" {
			myLogger.Debug("Failed init DB", zap.String("msg", err.Error()))
			return nil, err
		}
		repo = db
		break
	case config.FileStorage:
		fs, err := filestorage.InitFileStorageFile(cfg.FileStoragePath, myLogger)
		if err != nil {
			myLogger.Debug("Failed init filestorage", zap.String("msg", err.Error()))
			return nil, err
		}
		fmt.Println(fs)
		fmt.Println(err)
		repo = fs
		break
	case config.InMemory:
		inmemory, err := inmemorystorage.InitInMemoryStorage(map[string]string{}, myLogger)
		if err != nil {
			myLogger.Debug("Failed init inmemorystorage", zap.String("msg", err.Error()))
			return nil, err
		}
		repo = inmemory
		break
	}

	// fs := storage.InitFileStorageFile(cfg.FileStoragePath)
	// s := storage.InitStorage(map[string]string{}, fs, myLogger)

	// db, err := db.InitDB("localhost", "5432", "urls", "jf6y5SfnxsuR", "urls")
	srv, err := server.InitServer(cfg.URLPrefixRepo, repo, myLogger)
	// srv.MyDB = db
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
