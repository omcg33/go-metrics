package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMetric_ValidateStructError(t *testing.T) {
	rr := getMetric(&serviceMock{}, "unknown", "Alloc")

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetMetric_ValidationErrorsAs(t *testing.T) {
	rr := getMetric(&serviceMock{}, "gauge", "")

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetMetric_GaugeStatusOK(t *testing.T) {
	svc := &serviceMock{}
	svc.On("Gauge", "Alloc").Return(1.5, true)

	rr := getMetric(svc, "gauge", "Alloc")

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetMetric_CounterStatusOK(t *testing.T) {
	svc := &serviceMock{}
	svc.On("Counter", "Alloc").Return(1.5, true)

	rr := getMetric(svc, "counter", "Alloc")

	assert.Equal(t, http.StatusOK, rr.Code)
}

func getMetric(svc *serviceMock, metricType, name string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetPathValue("type", metricType)
	req.SetPathValue("name", name)

	NewController(svc).GetMetric(rr, req)

	return rr
}
