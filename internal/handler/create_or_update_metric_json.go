package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	models "github.com/omcg33/go-metrics/internal/model"
)

func (controller *Controller) CreateOrUpdateMetricJSON(res http.ResponseWriter, req *http.Request) {
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

	if data.MType != "gauge" && data.MType != "counter" {
		http.Error(res, "invalid metric type", http.StatusBadRequest)
		return
	}

	if data.ID == "" {
		http.Error(res, "metric name is required", http.StatusNotFound)
		return
	}

	switch data.MType {
	case models.Gauge:
		if data.Value == nil {
			http.Error(res, "metric value is required", http.StatusBadRequest)
			return
		}
		controller.service.CreateOrUpdateGauge(data.ID, *data.Value)
	case models.Counter:
		if data.Delta == nil {
			http.Error(res, "metric delta is required", http.StatusBadRequest)
			return
		}
		controller.service.CreateOrUpdateCounter(data.ID, *data.Delta)
	default:
		http.Error(res, "invalid metric type", http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write(jsonData)
}
