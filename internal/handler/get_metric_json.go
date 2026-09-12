package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

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

		res.Write([]byte(strconv.FormatFloat(value, 'f', -1, 64)))
	case models.Counter:
		value, isExist := controller.service.Counter(data.ID)
		if !isExist {
			http.Error(res, "metric name by name "+data.ID+"not found", http.StatusNotFound)
			return
		}

		res.Write([]byte(strconv.FormatInt(value, 10)))
	default:
		http.Error(res, "invalid metric type", http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
}
