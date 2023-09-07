package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	//nolint:depguard
	"github.com/kripsy/shortener/internal/client/clientcompress"
	//nolint:depguard
	"github.com/kripsy/shortener/internal/client/clientmodels"
	//nolint:depguard
	"github.com/kripsy/shortener/internal/client/clientutils"
)

func main() {
	endpoint := "http://localhost:8080/api/shorten"

	fmt.Println("Введите длинный URL")
	reader := bufio.NewReader(os.Stdin)
	long, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)

		return
	}
	long = strings.TrimSuffix(long, "\n")
	data := clientmodels.Requset{URL: long}

	//nolint:exhaustruct
	client := &http.Client{}

	reqData, err := json.Marshal(data)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)

		return
	}

	cb := clientcompress.Compress(string(reqData))
	//nolint:noctx
	request, err := http.NewRequest(http.MethodPost, endpoint, &cb)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)

		return
	}

	// set headers
	clientutils.SetHeaders(&request.Header)

	// отправляем запрос и получаем ответ
	response, err := client.Do(request)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)

		return
	}
	// выводим код ответа
	fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()
	// читаем поток из тела ответа
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)

		return
	}

	// и печатаем его
	fmt.Println(string(body))
	fmt.Println((clientcompress.Decompress(string(body))))
}
