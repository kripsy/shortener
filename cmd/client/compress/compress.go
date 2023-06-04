package compress

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"strings"
)

func Compress(data string) bytes.Buffer {
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

func Decompress(data string) string {
	rdata := strings.NewReader(data)
	r, err := gzip.NewReader(rdata)
	log.Println(r)
	if err != nil {
		log.Fatal(err)
	}
	s, _ := io.ReadAll(r)
	return (string(s))
}
