package middleware

import (
	"compress/gzip"
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
	//nolint:wrapcheck
	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode <= http.StatusIMUsed {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close closes gzip.Writer and send all data from buffer.
func (c *compressWriter) Close() error {
	//nolint:wrapcheck
	return c.zw.Close()
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
		//nolint:wrapcheck
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read return p bytes from compressReader.
func (c compressReader) Read(p []byte) (int, error) {
	//nolint:wrapcheck
	return c.zr.Read(p)
}

// Close closes Reader for compressReader.
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		//nolint:wrapcheck
		return err
	}
	//nolint:wrapcheck
	return c.zr.Close()
}
