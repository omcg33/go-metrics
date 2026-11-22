package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const decompressPayload = `{"id":"test","type":"gauge","value":1.5}`

func gzipBytes(t *testing.T, s string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write([]byte(s))
	require.NoError(t, err)
	require.NoError(t, zw.Close())
	return buf.Bytes()
}

// errReader всегда возвращает ошибку при чтении.
type errReader struct{}

func (errReader) Read([]byte) (int, error) { return 0, io.ErrUnexpectedEOF }

// serve прогоняет запрос через DecompressRequest и возвращает ответ,
// тело и длину, которые увидел next, а также признак того, что next был вызван.
func serve(t *testing.T, req *http.Request) (*httptest.ResponseRecorder, string, int64, bool) {
	t.Helper()
	var (
		gotBody string
		gotLen  int64
		called  bool
	)
	handler := DecompressRequest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		gotLen = r.ContentLength
		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		gotBody = string(b)
	}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec, gotBody, gotLen, called
}

func newGzipRequest(t *testing.T, encoding string, body []byte) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewReader(body))
	req.Header.Set("Content-Encoding", encoding)
	return req
}

// Правило: 1 тест — 1 assert.Equal.

func TestDecompressRequest_PassesBody(t *testing.T) {
	tests := []struct {
		name string
		req  func(t *testing.T) *http.Request
	}{
		{
			"WHEN Content-Encoding отсутствует MUST передать тело без изменений",
			func(t *testing.T) *http.Request {
				return httptest.NewRequest(http.MethodPost, "/update/", strings.NewReader(decompressPayload))
			},
		},
		{
			"WHEN Content-Encoding не gzip MUST передать тело без распаковки",
			func(t *testing.T) *http.Request {
				return newGzipRequest(t, "identity", []byte(decompressPayload))
			},
		},
		{
			"WHEN тело сжато gzip MUST передать распакованное тело",
			func(t *testing.T) *http.Request {
				return newGzipRequest(t, "gzip", gzipBytes(t, decompressPayload))
			},
		},
		{
			"WHEN Content-Encoding содержит gzip среди нескольких значений MUST передать распакованное тело",
			func(t *testing.T) *http.Request {
				return newGzipRequest(t, "deflate, gzip", gzipBytes(t, decompressPayload))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, body, _, _ := serve(t, tt.req(t))

			assert.Equal(t, decompressPayload, body)
		})
	}
}

func TestDecompressRequest_UpdatesContentLength(t *testing.T) {
	t.Run("WHEN тело сжато gzip MUST выставить ContentLength равным размеру распакованного тела", func(t *testing.T) {
		req := newGzipRequest(t, "gzip", gzipBytes(t, decompressPayload))

		_, _, gotLen, _ := serve(t, req)

		assert.Equal(t, int64(len(decompressPayload)), gotLen)
	})
}

func TestDecompressRequest_InvalidBody_StatusCode(t *testing.T) {
	valid := gzipBytes(t, decompressPayload)
	tests := []struct {
		name string
		body io.Reader
	}{
		{"WHEN gzip-тело пустое MUST вернуть 500", http.NoBody},
		{"WHEN тело не является валидным gzip MUST вернуть 500", strings.NewReader("not a gzip")},
		{"WHEN gzip-поток обрезан MUST вернуть 500", bytes.NewReader(valid[:len(valid)-8])},
		{"WHEN чтение тела завершается ошибкой MUST вернуть 500", errReader{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/update/", tt.body)
			req.Header.Set("Content-Encoding", "gzip")

			rec, _, _, _ := serve(t, req)

			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		})
	}
}

func TestDecompressRequest_InvalidBody_NextNotCalled(t *testing.T) {
	valid := gzipBytes(t, decompressPayload)
	tests := []struct {
		name string
		body io.Reader
	}{
		{"WHEN gzip-тело пустое MUST не вызывать next", http.NoBody},
		{"WHEN тело не является валидным gzip MUST не вызывать next", strings.NewReader("not a gzip")},
		{"WHEN gzip-поток обрезан MUST не вызывать next", bytes.NewReader(valid[:len(valid)-8])},
		{"WHEN чтение тела завершается ошибкой MUST не вызывать next", errReader{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/update/", tt.body)
			req.Header.Set("Content-Encoding", "gzip")

			_, _, _, called := serve(t, req)

			assert.Equal(t, false, called)
		})
	}
}
