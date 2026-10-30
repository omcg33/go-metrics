package middleware

import (
	"net/http"
	"strconv"
	"time"

	logs "github.com/omcg33/go-metrics/internal/logger"
	"go.uber.org/zap"
)

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

func (r *loggingResponseWriter) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}

func Logger(logger *logs.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()
			uri := r.RequestURI
			method := r.Method
			contentEncoding := r.Header.Get("Content-Encoding")

			responseData := &responseData{
				status: 0,
				size:   0,
			}

			lw := loggingResponseWriter{
				ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
				responseData:   responseData,
			}
			next.ServeHTTP(&lw, r)

			duration := time.Since(start)

			logger.Info(
				"Incoming HTTP request",
				zap.String("uri", uri),
				zap.String("method", method),
				zap.String("content-encoding", contentEncoding),
				zap.String("duration", duration.String()),
				zap.String("status", strconv.Itoa(responseData.status)),
				zap.String("size", strconv.Itoa(responseData.size)),
			)

		}
		return http.HandlerFunc(fn)
	}
}
