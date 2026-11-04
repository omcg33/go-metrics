package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type compressWriter struct {
	http.ResponseWriter
	compressor      *gzip.Writer
	isHeaderWritten bool
}

func shouldCompress(contentType string) bool {
	return strings.HasPrefix(contentType, "application/json") ||
		strings.HasPrefix(contentType, "text/html")
}

func (w *compressWriter) Close() error {
	if w.compressor != nil {
		return w.compressor.Close()
	}
	return nil
}

func (w *compressWriter) Write(b []byte) (int, error) {
	if !w.isHeaderWritten {
		w.WriteHeader(http.StatusOK)
	}

	if w.compressor != nil {
		return w.compressor.Write(b)
	}

	return w.ResponseWriter.Write(b)
}

func (w *compressWriter) WriteHeader(code int) {
	if w.isHeaderWritten {
		return
	}
	w.isHeaderWritten = true

	headers := w.Header()

	if shouldCompress(headers.Get("Content-Type")) {
		headers.Set("Content-Encoding", "gzip")
		headers.Del("Content-Length") // длина после сжатия изменится
		w.compressor = gzip.NewWriter(w.ResponseWriter)
	}

	w.ResponseWriter.WriteHeader(code)
}

func CompressResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		cw := &compressWriter{ResponseWriter: w}
		defer cw.Close()
		next.ServeHTTP(cw, r)
	})
}
