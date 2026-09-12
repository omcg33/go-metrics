package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMetricJSON_ValidateStructError(t *testing.T) {
	rr := getMetricJSON(NewMockMetricsService(t), `{"id":"Alloc","type":"unknown"}`)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetMetricJSON_ValidationErrorsAs(t *testing.T) {
	rr := getMetricJSON(NewMockMetricsService(t), `{"id":"","type":"gauge"}`)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetMetricJSON_GaugeStatusOK(t *testing.T) {
	svc := NewMockMetricsService(t)
	svc.On("Gauge", "Alloc").Return(1.5, true)

	rr := getMetricJSON(svc, `{"id":"Alloc","type":"gauge"}`)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "1.5", rr.Body.String())
}

func TestGetMetricJSON_GaugeOmitsTrailingZeros(t *testing.T) {
	svc := NewMockMetricsService(t)
	svc.On("Gauge", "testSetGet216").Return(65024.953, true)

	rr := getMetricJSON(svc, `{"id":"testSetGet216","type":"gauge"}`)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "65024.953", rr.Body.String())
}

func TestGetMetricJSON_CounterStatusOK(t *testing.T) {
	svc := NewMockMetricsService(t)
	svc.EXPECT().Counter("Alloc").Return(int64(1), true)

	rr := getMetricJSON(svc, `{"id":"Alloc","type":"counter"}`)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func getMetricJSON(svc *MockMetricsService, body string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	NewController(svc).GetMetricJSON(rr, req)

	return rr
}
