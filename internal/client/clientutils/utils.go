// Package clientutils provides general functionality for filling in the header for the client.
package clientutils

import "net/http"

func SetHeaders(h *http.Header) {
	// в заголовках запроса указываем кодировку
	h.Add("Content-Type", "application/json")
	h.Add("Accept-Encoding", "gzip")
	h.Add("Content-Encoding", "gzip")
}
