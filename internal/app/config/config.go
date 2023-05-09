package config

import (
	"flag"
)

type Config struct {
	// it's address for exec server
	UrlServer string

	// it's address for prefix in store short url
	UrlPrefixRepo string
}

func InitConfig() *Config {
	// декларируем наборы флагов для подкоманд
	urlServer := flag.String("a", "localhost:8080", "Enter address exec http server as ip_address:port")
	urlPrefixRepo := flag.String("b", "http://localhost:8080", "Enter address exec http server as ip_address:port")

	flag.Parse()

	c := Config{
		UrlServer:     *urlServer,
		UrlPrefixRepo: *urlPrefixRepo,
	}
	return &c
}
