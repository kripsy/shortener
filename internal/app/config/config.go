// Package config provides the functionality to collect the general configuration of the project.
// Config is designed to store URL server, URL prefix for repository, logging level, path for file storage, dsn for database.
package config

import (
	"flag"
	"fmt"
	"os"
)

// RepositoryType represent type of selected storage : inmemory, filestorage, postgresql.
type RepositoryType int

const (
	InMemory    RepositoryType = iota
	FileStorage RepositoryType = iota
	PostgresDB  RepositoryType = iota
)

// Config represent settings for service.
type Config struct {
	// URLServer is an address for exec server.
	URLServer string

	// URLPrefixRepo is an address for prefix in store short url.
	URLPrefixRepo string

	// LoggerLevel is a logger level.
	LoggerLevel string

	// FileStoragePath is a file storage path.
	FileStoragePath string

	// DatabaseDsn is a database conn string.
	DatabaseDsn string

	// RepositoryType is a field to check type of repo (db, file, rom memory).
	RepositoryType RepositoryType
}

// InitConfig return a pointer Config.
// Fields for config are taken from flags, or environment variables.
func InitConfig() *Config {
	var repositoryType RepositoryType

	// declare set of flags for subcommands.
	URLServer := flag.String("a", "localhost:8080", "Enter address exec http server as ip_address:port. Or use SERVER_ADDRESS env")
	URLPrefixRepo := flag.String("b", "http://localhost:8080", "Enter address exec http server as ip_address:port. Or use BASE_URL env")
	logLevel := flag.String("l", "Info", "log level: Debug, Info, Warn, Error and etc... Or use LOG_LEVEL env")
	fileStoragePath := flag.String("f", "", "set path for tmp file... Or use FILE_STORAGE_PATH env")
	databaseDsn := flag.String("d", "", "set path for database... Or use DATABASE_DSN env. Example host=localhost user=urls password=jf6y5SfnxsuR sslmode=disable port=5432")
	flag.Parse()

	if envSrvAddr := os.Getenv("SERVER_ADDRESS"); envSrvAddr != "" {
		*URLServer = envSrvAddr
	}

	if envPrefixAddr := os.Getenv("BASE_URL"); envPrefixAddr != "" {
		*URLPrefixRepo = envPrefixAddr
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		*logLevel = envLogLevel
	}

	if envDatabaseDsn := os.Getenv("DATABASE_DSN"); envDatabaseDsn != "" {
		*databaseDsn = envDatabaseDsn
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		*fileStoragePath = envFileStoragePath
	}

	switch {
	case *databaseDsn != "":
		fmt.Println("using PostgresDB")
		fmt.Println(*databaseDsn)
		if *databaseDsn == "postgres:5432/praktikum?sslmode=disable" {
			fmt.Println("change connstring")
			*databaseDsn = "postgres://postgres@localhost:5432/praktikum?sslmode=disable"

		}
		if *databaseDsn == "postgres:5432/postgres?sslmode=disable" {
			fmt.Println("change connstring")
			*databaseDsn = "postgres://postgres@localhost:5432/postgres?sslmode=disable"

		}
		repositoryType = PostgresDB

	case *fileStoragePath != "":
		fmt.Println("using FileStorage")
		repositoryType = FileStorage

	default:
		fmt.Println("using InMemory")
		repositoryType = InMemory

	}

	cfg := &Config{
		URLServer:       *URLServer,
		URLPrefixRepo:   *URLPrefixRepo,
		LoggerLevel:     *logLevel,
		FileStoragePath: *fileStoragePath,
		DatabaseDsn:     *databaseDsn,
		RepositoryType:  repositoryType,
	}
	return cfg
}
