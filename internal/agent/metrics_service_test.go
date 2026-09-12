package agent

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewService_NotNil(t *testing.T) {
	s := NewService("http://localhost:8080")
	t.Cleanup(func() { _ = s.Close() })

	assert.NotNil(t, s)
}

func TestNewService_ClientNotNil(t *testing.T) {
	s := NewService("http://localhost:8080")
	t.Cleanup(func() { _ = s.Close() })

	assert.NotNil(t, s.client)
}

func TestReport_GaugeMethod(t *testing.T) {
	req := reportAndCapture(t, Report{gauges: map[string]float64{"Alloc": 1.5}})

	assert.Equal(t, http.MethodPost, req.Method)
}

func TestReport_GaugePath(t *testing.T) {
	req := reportAndCapture(t, Report{gauges: map[string]float64{"Alloc": 1.5}})

	assert.Equal(t, "/update/", req.Path)
}

func TestReport_GaugeContentType(t *testing.T) {
	req := reportAndCapture(t, Report{gauges: map[string]float64{"Alloc": 1.5}})

	assert.Equal(t, "application/json", req.ContentType)
}

func TestReport_GaugeBody(t *testing.T) {
	req := reportAndCapture(t, Report{gauges: map[string]float64{"Alloc": 1.5}})

	assert.JSONEq(t, `{"id":"Alloc","type":"gauge","value":1.5}`, string(req.Body))
}

func TestReport_CounterMethod(t *testing.T) {
	req := reportAndCapture(t, Report{counters: map[string]int64{"PollCount": 42}})

	assert.Equal(t, http.MethodPost, req.Method)
}

func TestReport_CounterPath(t *testing.T) {
	req := reportAndCapture(t, Report{counters: map[string]int64{"PollCount": 42}})

	assert.Equal(t, "/update/", req.Path)
}

func TestReport_CounterContentType(t *testing.T) {
	req := reportAndCapture(t, Report{counters: map[string]int64{"PollCount": 42}})

	assert.Equal(t, "application/json", req.ContentType)
}

func TestReport_CounterBody(t *testing.T) {
	req := reportAndCapture(t, Report{counters: map[string]int64{"PollCount": 42}})

	assert.JSONEq(t, `{"id":"PollCount","type":"counter","delta":42}`, string(req.Body))
}

func TestReport_EmptyMakesNoRequests(t *testing.T) {
	var n int
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		n++
	}))
	t.Cleanup(srv.Close)

	s := NewService(srv.URL)
	t.Cleanup(func() { _ = s.Close() })
	s.Report(Report{})

	assert.Equal(t, 0, n)
}

func TestReport_DoesNotPanicWhenGaugePostFails(t *testing.T) {
	s := NewService(closedServerURL(t))
	t.Cleanup(func() { _ = s.Close() })

	assert.NotPanics(t, func() {
		s.Report(Report{gauges: map[string]float64{"Alloc": 1}})
	})
}

func TestReport_DoesNotPanicWhenCounterPostFails(t *testing.T) {
	s := NewService(closedServerURL(t))
	t.Cleanup(func() { _ = s.Close() })

	assert.NotPanics(t, func() {
		s.Report(Report{counters: map[string]int64{"PollCount": 1}})
	})
}

type capturedRequest struct {
	Method      string
	Path        string
	ContentType string
	Body        []byte
}

func reportAndCapture(t *testing.T, report Report) capturedRequest {
	t.Helper()

	var got capturedRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		got = capturedRequest{
			Method:      r.Method,
			Path:        r.URL.Path,
			ContentType: r.Header.Get("Content-Type"),
			Body:        body,
		}
	}))
	t.Cleanup(srv.Close)

	s := NewService(srv.URL)
	t.Cleanup(func() { _ = s.Close() })
	s.Report(report)

	return got
}

func closedServerURL(t *testing.T) string {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	url := srv.URL
	srv.Close()

	return url
}
