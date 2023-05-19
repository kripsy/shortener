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
}

func InitConfig() *Config {
	// декларируем наборы флагов для подкоманд
	URLServer := flag.String("a", "localhost:8080", "Enter address exec http server as ip_address:port. Or use SERVER_ADDRESS env")
	URLPrefixRepo := flag.String("b", "http://localhost:8080", "Enter address exec http server as ip_address:port. Or use BASE_URL env")
	logLevel := flag.String("l", "Info", "log level: Debug, Info, Warn, Error and etc... Or use LOG_LEVEL env")

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

	cfg := &Config{
		URLServer:     *URLServer,
		URLPrefixRepo: *URLPrefixRepo,
		LoggerLevel:   *logLevel,
	}
	return cfg
}
