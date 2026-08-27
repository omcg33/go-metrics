package agent

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewService_NotNil(t *testing.T) {
	assert.NotNil(t, NewService("http://localhost:8080"))
}

func TestNewService_ClientNotNil(t *testing.T) {
	assert.NotNil(t, NewService("http://localhost:8080").client)
}

func TestReport_GaugeMethod(t *testing.T) {
	req := reportAndCapture(t, Report{gauges: map[string]float64{"Alloc": 1.5}})

	assert.Equal(t, http.MethodPost, req.Method)
}

func TestReport_GaugePath(t *testing.T) {
	req := reportAndCapture(t, Report{gauges: map[string]float64{"Alloc": 1.5}})

	assert.Equal(t, "/update/gauge/Alloc/1.500000", req.URL.Path)
}

func TestReport_GaugeContentType(t *testing.T) {
	req := reportAndCapture(t, Report{gauges: map[string]float64{"Alloc": 1.5}})

	assert.Equal(t, "text/plain", req.Header.Get("Content-Type"))
}

func TestReport_CounterMethod(t *testing.T) {
	req := reportAndCapture(t, Report{counters: map[string]int64{"PollCount": 42}})

	assert.Equal(t, http.MethodPost, req.Method)
}

func TestReport_CounterPath(t *testing.T) {
	req := reportAndCapture(t, Report{counters: map[string]int64{"PollCount": 42}})

	assert.Equal(t, "/update/counter/PollCount/42", req.URL.Path)
}

func TestReport_CounterContentType(t *testing.T) {
	req := reportAndCapture(t, Report{counters: map[string]int64{"PollCount": 42}})

	assert.Equal(t, "text/plain", req.Header.Get("Content-Type"))
}

func TestReport_EmptyMakesNoRequests(t *testing.T) {
	var n int
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		n++
	}))
	t.Cleanup(srv.Close)

	s := NewService(srv.URL)
	s.Report(Report{})

	assert.Equal(t, 0, n)
}

func TestReport_PanicsWhenGaugePostFails(t *testing.T) {
	s := NewService(closedServerURL(t))

	assert.Panics(t, func() {
		s.Report(Report{gauges: map[string]float64{"Alloc": 1}})
	})
}

func TestReport_PanicsWhenCounterPostFails(t *testing.T) {
	s := NewService(closedServerURL(t))

	assert.Panics(t, func() {
		s.Report(Report{counters: map[string]int64{"PollCount": 1}})
	})
}

func reportAndCapture(t *testing.T, report Report) *http.Request {
	t.Helper()

	var got *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Clone(r.Context())
	}))
	t.Cleanup(srv.Close)

	s := NewService(srv.URL)
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
