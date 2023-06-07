package config

import (
	"flag"
	"os"
)

type Config struct {
	// it's address for exec server
	URLServer string

	// it's address for prefix in store short url
	URLPrefixRepo string

	// it's logger level
	LoggerLevel string

	// it's file storage path
	FileStoragePath string

	// it's database conn string
	DatabaseDsn string
}

func InitConfig() *Config {
	// декларируем наборы флагов для подкоманд
	URLServer := flag.String("a", "localhost:8080", "Enter address exec http server as ip_address:port. Or use SERVER_ADDRESS env")
	URLPrefixRepo := flag.String("b", "http://localhost:8080", "Enter address exec http server as ip_address:port. Or use BASE_URL env")
	logLevel := flag.String("l", "Info", "log level: Debug, Info, Warn, Error and etc... Or use LOG_LEVEL env")
	fileStoragePath := flag.String("f", "/tmp/short-url-db.json", "set path for tmp file... Or use FILE_STORAGE_PATH env")
	databaseDsn := flag.String("d", "", "set path for database... Or use DATABASE_DSN env. Example host=localhost user=urls password=jf6y5SfnxsuR sslmode=disable port=5432")
	//host=localhost user=urls password=jf6y5SfnxsuR sslmode=disable port=5432
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

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		*fileStoragePath = envFileStoragePath
	}

	if envDatabaseDsn := os.Getenv("DATABASE_DSN"); envDatabaseDsn != "" {
		*databaseDsn = envDatabaseDsn
	}

	cfg := &Config{
		URLServer:       *URLServer,
		URLPrefixRepo:   *URLPrefixRepo,
		LoggerLevel:     *logLevel,
		FileStoragePath: *fileStoragePath,
		DatabaseDsn:     *databaseDsn,
	}
	return cfg
}
