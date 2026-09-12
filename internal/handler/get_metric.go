package handler

import (
	"net/http"
	"strconv"

	models "github.com/omcg33/go-metrics/internal/model"
)

type GetMetricParams struct {
	Type string `validate:"required,metric_type"`
	Name string `validate:"required"`
}

func (controller *Controller) GetMetric(res http.ResponseWriter, req *http.Request) {
	metricType := req.PathValue("type")
	metricName := req.PathValue("name")

	if metricType != "gauge" || metricType != "counter" {
		http.Error(res, "invalid metric type", http.StatusBadRequest)
	}

	if metricName == "" {
		http.Error(res, "metric name is required", http.StatusNotFound)
	}

	switch metricType {
	case models.Gauge:
		value, isExist := controller.service.Gauge(metricName)
		if !isExist {
			http.Error(res, "metric name by name "+metricName+" not found", http.StatusNotFound)
			return
		}

		res.Write([]byte(strconv.FormatFloat(value, 'f', -1, 64)))
	case models.Counter:
		value, isExist := controller.service.Counter(metricName)
		if !isExist {
			http.Error(res, "metric name by name "+metricName+"not found", http.StatusNotFound)
			return
		}

		res.Write([]byte(strconv.FormatInt(value, 10)))
	default:
		http.Error(res, "invalid metric type", http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
}
