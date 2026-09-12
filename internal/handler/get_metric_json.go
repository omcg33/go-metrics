package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	models "github.com/omcg33/go-metrics/internal/model"
)

func (controller *Controller) GetMetricJSON(res http.ResponseWriter, req *http.Request) {
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
		value, isExist := controller.service.Gauge(data.ID)

		if !isExist {
			http.Error(res, "metric name by name "+data.ID+" not found", http.StatusNotFound)
			return
		}

		jsonData, err := json.Marshal(models.Metrics{
			ID:    data.ID,
			MType: data.MType,
			Value: &value,
		})

		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(jsonData)
	case models.Counter:
		value, isExist := controller.service.Counter(data.ID)

		if !isExist {
			http.Error(res, "metric name by name "+data.ID+"not found", http.StatusNotFound)
			return
		}

		jsonData, err := json.Marshal(models.Metrics{
			ID:    data.ID,
			MType: data.MType,
			Delta: &value,
		})

		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(jsonData)
	default:
		http.Error(res, "invalid metric type", http.StatusBadRequest)
		return
	}
}
