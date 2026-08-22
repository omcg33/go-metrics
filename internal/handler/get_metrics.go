package handler

import (
	"encoding/json"
	"net/http"
)

func (controller *Controller) GetMetrics(res http.ResponseWriter, req *http.Request) {
	metrics := controller.service.Metrics()
	json.NewEncoder(res).Encode(metrics)
}