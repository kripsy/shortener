package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/kripsy/shortener/internal/client/clientcompress"
	"github.com/kripsy/shortener/internal/client/clientmodels"
	"github.com/kripsy/shortener/internal/client/clientutils"
)

func main() {

	endpoint := "http://localhost:8080/api/shorten"
	// контейнер данных для запроса
	data := clientmodels.Requset{}
	// приглашение в консоли
	fmt.Println("Введите длинный URL")
	// открываем потоковое чтение из консоли
	reader := bufio.NewReader(os.Stdin)
	// читаем строку из консоли
	long, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return
	}
	long = strings.TrimSuffix(long, "\n")
	// заполняем контейнер данными

	data.URL = long

	// добавляем HTTP-клиент
	client := &http.Client{}
	// пишем запрос
	// запрос методом POST должен, помимо заголовков, содержать тело
	// тело должно быть источником потокового чтения io.Reader
	reqData, err := json.Marshal(data)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return
	}

	cb := clientcompress.Compress(string(reqData))

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
