package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type compressWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w compressWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func CompressResponse(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(compressWriter{ResponseWriter: w, Writer: gz}, r)
	}

	return http.HandlerFunc(fn)
}
