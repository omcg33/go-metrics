package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMetric_ValidateStructError(t *testing.T) {
	rr := getMetric(NewMockMetricsService(t), "unknown", "Alloc")

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetMetric_ValidationErrorsAs(t *testing.T) {
	rr := getMetric(NewMockMetricsService(t), "gauge", "")

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetMetric_GaugeStatusOK(t *testing.T) {
	svc := NewMockMetricsService(t)
	svc.On("Gauge", "Alloc").Return(1.5, true)

	rr := getMetric(svc, "gauge", "Alloc")

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "1.5", rr.Body.String())
}

func TestGetMetric_GaugeOmitsTrailingZeros(t *testing.T) {
	svc := NewMockMetricsService(t)
	svc.On("Gauge", "testSetGet216").Return(65024.953, true)

	rr := getMetric(svc, "gauge", "testSetGet216")

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "65024.953", rr.Body.String())
}

func TestGetMetric_CounterStatusOK(t *testing.T) {
	svc := NewMockMetricsService(t)
	svc.EXPECT().Counter("Alloc").Return(int64(1), true)

	rr := getMetric(svc, "counter", "Alloc")

	assert.Equal(t, http.StatusOK, rr.Code)
}

func getMetric(svc *MockMetricsService, metricType, name string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetPathValue("type", metricType)
	req.SetPathValue("name", name)

	NewController(svc).GetMetric(rr, req)

	return rr
}
