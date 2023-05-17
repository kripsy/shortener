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
}

func InitConfig() *Config {
	// декларируем наборы флагов для подкоманд
	URLServer := flag.String("a", "localhost:8080", "Enter address exec http server as ip_address:port")
	URLPrefixRepo := flag.String("b", "http://localhost:8080", "Enter address exec http server as ip_address:port")

	flag.Parse()

	if envSrvAddr := os.Getenv("SERVER_ADDRESS"); envSrvAddr != "" {

		*URLServer = envSrvAddr
	}
	if envPrefixAddr := os.Getenv("BASE_URL"); envPrefixAddr != "" {
		*URLPrefixRepo = envPrefixAddr
	}

	cfg := &Config{
		URLServer:     *URLServer,
		URLPrefixRepo: *URLPrefixRepo,
	}
	return cfg
}
