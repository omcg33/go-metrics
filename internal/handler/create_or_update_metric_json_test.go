package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrUpdateMetricJson_ValidateStructError(t *testing.T) {
	rr := createOrUpdateJSON(NewMockMetricsService(t), `{"id":"Alloc","type":"unknown","value":1.5}`)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateOrUpdateMetricJson_ValidationErrorsAs(t *testing.T) {
	rr := createOrUpdateJSON(NewMockMetricsService(t), `{"id":"","type":"gauge","value":1.5}`)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestCreateOrUpdateMetricJson_ParseFloatError(t *testing.T) {
	rr := createOrUpdateJSON(NewMockMetricsService(t), `{"id":"Alloc","type":"gauge","value":"abc"}`)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateOrUpdateMetricJson_ParseIntError(t *testing.T) {
	rr := createOrUpdateJSON(NewMockMetricsService(t), `{"id":"PollCount","type":"counter","delta":1.5}`)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateOrUpdateMetricJson_StatusOK(t *testing.T) {
	svc := NewMockMetricsService(t)
	svc.EXPECT().CreateOrUpdateGauge("Alloc", 1.5)

	rr := createOrUpdateJSON(svc, `{"id":"Alloc","type":"gauge","value":1.5}`)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func createOrUpdateJSON(svc *MockMetricsService, body string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	NewController(svc).CreateOrUpdateMetricJson(rr, req)

	return rr
}
