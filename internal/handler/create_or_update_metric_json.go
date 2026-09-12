package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	models "github.com/omcg33/go-metrics/internal/model"
)

func (controller *Controller) CreateOrUpdateMetricJson(res http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	var data models.Metrics

	_, err := buf.ReadFrom(req.Body)

	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &data); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if data.MType != "gauge" || data.MType != "counter" {
		http.Error(res, "invalid metric type", http.StatusBadRequest)
	}

	if data.ID == "" {
		http.Error(res, "metric name is required", http.StatusNotFound)
	}

	switch data.MType {
	case models.Gauge:
		controller.service.CreateOrUpdateGauge(data.ID, *data.Value)
	case models.Counter:
		controller.service.CreateOrUpdateCounter(data.ID, *data.Delta)
	default:
		http.Error(res, "invalid metric type", http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)

	res.WriteHeader(http.StatusOK)
}
