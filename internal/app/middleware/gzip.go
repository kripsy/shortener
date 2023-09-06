package middleware

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
)

// compressWriter realize interface http.ResponseWriter,
// allows compress data transparent to the server, set correct headers.
type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// newCompressWriter return new compressWriter pointer.
func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header return Header of writer.
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Write call method Write of gzip writer.
func (c *compressWriter) Write(p []byte) (int, error) {
	n, err := c.zw.Write(p)
	if err != nil {
		return n, fmt.Errorf("%w", err)
	}

	return n, nil
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode <= http.StatusIMUsed {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close closes gzip.Writer and send all data from buffer.
func (c *compressWriter) Close() error {
	if err := c.zw.Close(); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// compressReader implements the io.ReadCloser interface and makes it transparent to the server
// decompress data received from the client.
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// newCompressReader return new compressReader pointer and error if exists.
func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read return p bytes from compressReader.
//
//nolint:nonamedreturns
func (c compressReader) Read(p []byte) (n int, err error) {
	n, err = c.zr.Read(p)
	if err != nil {
		return n, fmt.Errorf("%w", err)
	}

	return n, nil
}

// Close closes Reader for compressReader.
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
