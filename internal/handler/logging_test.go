package handler

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggingLogsMethodAndURL(t *testing.T) {
	var buf bytes.Buffer
	origWriter := log.Writer()
	origFlags := log.Flags()
	log.SetOutput(&buf)
	log.SetFlags(0)
	t.Cleanup(func() {
		log.SetOutput(origWriter)
		log.SetFlags(origFlags)
	})

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/update/gauge/Alloc/1.23", nil)
	rr := httptest.NewRecorder()

	Logging(next).ServeHTTP(rr, req)

	got := buf.String()
	if !strings.Contains(got, http.MethodPost) {
		t.Fatalf("log %q does not contain request method", got)
	}
	if !strings.Contains(got, "/update/gauge/Alloc/1.23") {
		t.Fatalf("log %q does not contain request URL", got)
	}
}
