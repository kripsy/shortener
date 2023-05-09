package config

type Config struct {
	Url  string
	Port string
}

func InitConfig(url, port string) *Config {
	c := Config{
		Url:  url,
		Port: port,
	}
	return &c
}
