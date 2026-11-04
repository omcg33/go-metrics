package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompressResponse(t *testing.T) {
	const body = `{"id":"test","type":"gauge","value":1.5}`

	tests := []struct {
		name           string
		acceptEncoding string
		contentType    string
		explicitHeader bool // вызывать WriteHeader явно перед Write
		wantGzip       bool
	}{
		{"json + gzip, только Write", "gzip", "application/json", false, true},
		{"json с charset + gzip", "gzip", "application/json; charset=utf-8", false, true},
		{"html + gzip", "gzip", "text/html", false, true},
		{"json + gzip, явный WriteHeader", "gzip", "application/json", true, true},
		{"text/plain не сжимается", "gzip", "text/plain", false, false},
		{"без Content-Type не сжимается", "gzip", "", false, false},
		{"клиент не поддерживает gzip", "", "application/json", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := CompressResponse(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.contentType != "" {
					w.Header().Set("Content-Type", tt.contentType)
				}
				if tt.explicitHeader {
					w.WriteHeader(http.StatusOK)
				}
				_, err := w.Write([]byte(body))
				require.NoError(t, err)
			}))

			req := httptest.NewRequest(http.MethodPost, "/update/", nil)
			if tt.acceptEncoding != "" {
				req.Header.Set("Accept-Encoding", tt.acceptEncoding)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()
			assert.Equal(t, http.StatusOK, res.StatusCode)

			if !tt.wantGzip {
				assert.Empty(t, res.Header.Get("Content-Encoding"))
				got, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Equal(t, body, string(got))
				return
			}

			assert.Equal(t, "gzip", res.Header.Get("Content-Encoding"))
			zr, err := gzip.NewReader(res.Body)
			require.NoError(t, err)
			got, err := io.ReadAll(zr)
			require.NoError(t, err)
			assert.Equal(t, body, string(got))
		})
	}
}

func TestCompressResponse_RemovesContentLength(t *testing.T) {
	handler := CompressResponse(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", "100")
		_, _ = w.Write([]byte(`{}`))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Empty(t, rec.Result().Header.Get("Content-Length"))
}

func TestCompressResponse_KeepsStatusCode(t *testing.T) {
	handler := CompressResponse(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("not found"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, "gzip", rec.Header().Get("Content-Encoding"))
}
