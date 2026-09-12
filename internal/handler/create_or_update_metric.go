package handler

import (
	"net/http"
	"strconv"

	models "github.com/omcg33/go-metrics/internal/model"
)

func (controller *Controller) CreateOrUpdateMetric(res http.ResponseWriter, req *http.Request) {

	metricType := req.PathValue("type")
	metricName := req.PathValue("name")
	metricValue := req.PathValue("value")

	if metricType != "gauge" && metricType != "counter" {
		http.Error(res, "invalid metric type", http.StatusBadRequest)
	}

	if metricName == "" {
		http.Error(res, "metric name is required", http.StatusNotFound)
	}

	switch metricType {
	case models.Gauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(res, "invalid metric value", http.StatusBadRequest)
			return
		}
		controller.service.CreateOrUpdateGauge(metricName, value)
	case models.Counter:
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(res, "invalid metric value", http.StatusBadRequest)
			return
		}
		controller.service.CreateOrUpdateCounter(metricName, value)
	default:
		http.Error(res, "invalid metric type", http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
}
