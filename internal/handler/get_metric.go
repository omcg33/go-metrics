package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	models "github.com/omcg33/go-metrics/internal/model"
)

type GetMetricParams struct {
	Type string `validate:"required,metric_type"`
	Name string `validate:"required"`
}

func (controller *Controller) GetMetric(res http.ResponseWriter, req *http.Request) {
	params := GetMetricParams{
		Type: req.PathValue("type"),
		Name: req.PathValue("name"),
	}

	if err := controller.validator.Struct(params); err != nil {
		var validationErrors validator.ValidationErrors

		if errors.As(err, &validationErrors) {
			for _, fieldErr := range validationErrors {
				switch fieldErr.Field() {
				case "Type":
					http.Error(res, "invalid metric type", http.StatusBadRequest)
				case "Name":
					http.Error(res, "metric name is required", http.StatusNotFound)
				}
			}
			return
		}
	}

	switch params.Type {
	case models.Gauge:
		value, isExist := controller.service.Gauge(params.Name)
		if !isExist {
			http.Error(res, "metric name by name "+params.Name+" not found", http.StatusNotFound)
			return
		}

		res.Write([]byte(strconv.FormatFloat(value, 'f', -1, 64)))
	case models.Counter:
		value, isExist := controller.service.Counter(params.Name)
		if !isExist {
			http.Error(res, "metric name by name "+params.Name+"not found", http.StatusNotFound)
			return
		}

		res.Write([]byte(strconv.FormatInt(value, 10)))
	default:
		http.Error(res, "invalid metric type", http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
}
