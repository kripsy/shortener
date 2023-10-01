// Package config provides the functionality to collect the general configuration of the project.
// Config is designed to store URL server,
// URL prefix for repository, logging level, path for file storage, dsn for database.
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
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
	URLServer string `json:"server_address,omitempty"`

	// URLPrefixRepo is an address for prefix in store short url.
	URLPrefixRepo string `json:"base_url,omitempty"`

	// LoggerLevel is a logger level.
	LoggerLevel string

	// FileStoragePath is a file storage path.
	FileStoragePath string `json:"file_storage_path,omitempty"`

	// DatabaseDsn is a database conn string.
	DatabaseDsn string `json:"database_dsn,omitempty"`

	// RepositoryType is a field to check type of repo (db, file, rom memory).
	RepositoryType RepositoryType

	// EnableHTTPS is a field to check is tls encryption
	EnableHTTPS string `json:"enable_https,omitempty"`

	// ConfigFilePath is a field with path to the config file
	ConfigFilePath string
}

// InitConfig return a pointer Config.
// Fields for config are taken from flags, or environment variables.
//
//nolint:cyclop
func InitConfig() *Config {
	var repositoryType RepositoryType

	// declare set of flags for subcommands.
	URLServer := flag.String("a",
		"localhost:8080",
		"Enter address exec http server as ip_address:port. Or use SERVER_ADDRESS env")
	URLPrefixRepo := flag.String("b",
		"http://localhost:8080",
		"Enter address exec http server as ip_address:port. Or use BASE_URL env")
	logLevel := flag.String("l", "Info", "log level: Debug, Info, Warn, Error and etc... Or use LOG_LEVEL env")
	fileStoragePath := flag.String("f", "", "set path for tmp file... Or use FILE_STORAGE_PATH env")
	databaseDsn := flag.String("d",
		"",
		`set path for database... Or use DATABASE_DSN env. 
		Example host=localhost user=urls password=jf6y5SfnxsuR sslmode=disable port=5432`)

	enableHTTPS := flag.String("s", "", "set tls encryption... Or use ENABLE_HTTPS env")

	configFilePath := flag.String("c", "", "set filepath for config... Or use CONFIG env")

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

	if envEnableHTTPS := os.Getenv("ENABLE_HTTPS"); envEnableHTTPS != "" {
		*enableHTTPS = envEnableHTTPS
	}

	if envConfigFilePath := os.Getenv("CONFIG"); envConfigFilePath != "" {
		*configFilePath = envConfigFilePath
	}

	if *configFilePath != "" {
		tempURLServer,
			tempURLPrefixRepo,
			tempDatabaseDsn,
			tempFileStoragePath,
			tempEnableHTTPS,
			err := updateConfigAttrFromFile(*configFilePath,
			*URLServer,
			*URLPrefixRepo,
			*databaseDsn,
			*fileStoragePath,
			*enableHTTPS)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			*URLServer = tempURLServer
			*URLPrefixRepo = tempURLPrefixRepo
			*databaseDsn = tempDatabaseDsn
			*fileStoragePath = tempFileStoragePath
			*enableHTTPS = tempEnableHTTPS
		}
	}

	repositoryType = setRepositoryType(*databaseDsn, *fileStoragePath)

	if repositoryType == PostgresDB {
		if *databaseDsn == "postgres:5432/praktikum?sslmode=disable" {
			fmt.Println("change connstring")
			*databaseDsn = "postgres://postgres@localhost:5432/praktikum?sslmode=disable"
		}
		if *databaseDsn == "postgres:5432/postgres?sslmode=disable" {
			fmt.Println("change connstring")
			*databaseDsn = "postgres://postgres@localhost:5432/postgres?sslmode=disable"
		}
	}
	*URLPrefixRepo = setHTTPS(*URLPrefixRepo, *enableHTTPS)

	return &Config{
		URLServer:       *URLServer,
		URLPrefixRepo:   *URLPrefixRepo,
		LoggerLevel:     *logLevel,
		FileStoragePath: *fileStoragePath,
		DatabaseDsn:     *databaseDsn,
		EnableHTTPS:     *enableHTTPS,
		RepositoryType:  repositoryType,
		ConfigFilePath:  *configFilePath,
	}
}

func updateConfigAttrFromFile(path string,
	urlServer,
	urlPrefixRepo,
	databaseDsn,
	fileStoragePath,
	enableHTTPS string) (string,
	string,
	string,
	string,
	string,
	error) {
	fmt.Println("read config file")

	inputParams := map[string]string{
		"server_address":    urlServer,
		"base_url":          urlPrefixRepo,
		"database_dsn":      databaseDsn,
		"file_storage_path": fileStoragePath,
		"enable_https":      enableHTTPS,
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("err read config file %v ", err)

		return "", "", "", "", "", fmt.Errorf("%w", err)
	}
	cfg := map[string]interface{}{}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		fmt.Printf("error unmarshall config %v: ", err)

		return "", "", "", "", "", fmt.Errorf("error unmarshall config %w", err)
	}

	for k, v := range inputParams {
		//nolint:nestif
		if v == "" {
			if k == "enable_https" {
				val, ok := cfg[k].(bool)
				if ok && val {
					inputParams[k] = "true"
				} else {
					inputParams[k] = "false"
				}
			} else {
				val, ok := cfg[k].(string)
				if ok {
					inputParams[k] = val
				}
			}
		}
	}

	return inputParams["server_address"],
		inputParams["base_url"],
		inputParams["database_dsn"],
		inputParams["file_storage_path"],
		inputParams["enable_https"],
		nil
}

func setRepositoryType(dsn, filePath string) RepositoryType {
	switch {
	case dsn != "":

		return PostgresDB
	case filePath != "":
		fmt.Println("using FileStorage")

		return FileStorage
	default:
		fmt.Println("using InMemory")

		return InMemory
	}
}

func setHTTPS(urlPrefixRepo, enableHTTPS string) string {
	if enableHTTPS != "" {
		fmt.Println("SET HTTPS")

		return strings.Replace(urlPrefixRepo, "http", "https", 1)
	}

	return urlPrefixRepo
}
