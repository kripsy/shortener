package config

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type Config struct {
	Url  string
	Port string
}

func InitConfig(url, port string) *Config {
	// декларируем наборы флагов для подкоманд
	urlServer := flag.String("a", "localhost:8080", "Enter address exec http server as ip_address:port")

	flag.Parse()
	a, b := parseUrl(*urlServer)

	fmt.Println(a)
	fmt.Println(b)
	// flag.Func("a", "ip address for exec server", func(flagValue string) error {

	// 	url = strings.Split(flagValue, ":")[0]

	// 	if strings.Split(flagValue, ":")[0] == "" {
	// 		return fmt.Errorf("ip address is empty")
	// 	}
	// 	// try parse port to int
	// 	_, err := strconv.Atoi(strings.Split(flagValue, ":")[1])
	// 	if err != nil {
	// 		// ... handle error
	// 		return err
	// 	}
	// 	port = strings.Split(flagValue, ":")[1]
	// 	return nil
	// })

	// fmt.Println(url)
	// fmt.Println(port)
	c := Config{
		Url:  url,
		Port: port,
	}
	return &c
}

func parseUrl(str string) (string, string) {
	url := strings.Split(str, ":")[0]
	var port string

	if len(strings.Split(str, ":")) < 2 {
		url = "localhost"
		port = "8080"
		fmt.Println("incorrect flag a, set default server address")
		return url, port
	}

	if url == "" {
		url = "localhost"
	}
	// try parse port to int
	_, err := strconv.Atoi(strings.Split(str, ":")[1])
	if err != nil {
		fmt.Println("incorrect flag a, set default server port")
		port = "8080"
		return url, port
	}
	port = strings.Split(str, ":")[1]
	return url, port
}
