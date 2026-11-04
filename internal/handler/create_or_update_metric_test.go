package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrUpdateMetric_ValidateStructError(t *testing.T) {
	rr := createOrUpdate(NewMockMetricsService(t), "unknown", "Alloc", "1.5")

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateOrUpdateMetric_ValidationErrorsAs(t *testing.T) {
	rr := createOrUpdate(NewMockMetricsService(t), "gauge", "", "1.5")

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestCreateOrUpdateMetric_ParseFloatError(t *testing.T) {
	rr := createOrUpdate(NewMockMetricsService(t), "gauge", "Alloc", "abc")

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateOrUpdateMetric_ParseIntError(t *testing.T) {
	rr := createOrUpdate(NewMockMetricsService(t), "counter", "PollCount", "1.5")

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateOrUpdateMetric_StatusOK(t *testing.T) {
	svc := NewMockMetricsService(t)
	svc.EXPECT().CreateOrUpdateGauge("Alloc", 1.5)

	rr := createOrUpdate(svc, "gauge", "Alloc", "1.5")

	assert.Equal(t, http.StatusOK, rr.Code)
}

func createOrUpdate(svc *MockMetricsService, metricType, name, value string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.SetPathValue("type", metricType)
	req.SetPathValue("name", name)
	req.SetPathValue("value", value)

	NewController(svc).CreateOrUpdateMetric(rr, req)

	return rr
}
