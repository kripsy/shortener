package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {

	endpoint := "http://localhost:8080/api/shorten"
	// контейнер данных для запроса
	data := url.Values{}
	// приглашение в консоли
	fmt.Println("Введите длинный URL")
	// открываем потоковое чтение из консоли
	reader := bufio.NewReader(os.Stdin)
	// читаем строку из консоли
	long, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	long = strings.TrimSuffix(long, "\n")
	// заполняем контейнер данными
	data.Set("url", long)
	// добавляем HTTP-клиент
	client := &http.Client{}
	// пишем запрос
	// запрос методом POST должен, помимо заголовков, содержать тело
	// тело должно быть источником потокового чтения io.Reader
	a := `{"url":"https://ya.ru"}`

	cb := compress(a)

	request, err := http.NewRequest(http.MethodPost, endpoint, &cb)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	// в заголовках запроса указываем кодировку
	request.Header.Add("Content-Type", "application/json")

	request.Header.Add("Accept-Encoding", "gzip")
	request.Header.Add("Content-Encoding", "gzip")

	// отправляем запрос и получаем ответ
	response, err := client.Do(request)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	// выводим код ответа
	fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()
	// читаем поток из тела ответа
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// и печатаем его
	fmt.Println(string(body))
	fmt.Println((decompress(string(body))))
}

func compress(data string) bytes.Buffer {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write([]byte(data))
	if err != nil {
		log.Fatal(err)
	}

	if err := zw.Close(); err != nil {
		log.Fatal(err)
	}
	return buf
}

func decompress(data string) string {
	rdata := strings.NewReader(data)
	r, err := gzip.NewReader(rdata)
	log.Println(r)
	if err != nil {
		log.Fatal(err)
	}
	s, _ := ioutil.ReadAll(r)
	return (string(s))
}
