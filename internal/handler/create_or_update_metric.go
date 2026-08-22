package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"

	models "github.com/omcg33/go-metrics/internal/model"
)

type MetricType string

type Params struct {
	Type MetricType `validate:"required,metric_type"`
	Name string     `validate:"required"`
	Value string     `validate:"required"`
}

var validate = newValidate()

func newValidate() *validator.Validate {
	v := validator.New()
	v.RegisterAlias("metric_type", fmt.Sprintf("oneof=%s %s", models.Gauge, models.Counter))
	return v
}

func (controller *Controller) CreateOrUpdateMetric(res http.ResponseWriter, req *http.Request) {
	params := Params{
		Type: MetricType(req.PathValue("type")),
		Name: req.PathValue("name"),
		Value: req.PathValue("value"),
	}

	if err := validate.Struct(params); err != nil {

		var validationErrors validator.ValidationErrors

		if errors.As(err, &validationErrors) {
			for _, fieldErr := range validationErrors {
				switch fieldErr.Field() {
				case "Type", "Value":
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
		value, err := strconv.ParseFloat(params.Value, 64)
		if err != nil {
			http.Error(res, "invalid metric value", http.StatusBadRequest)
			return
		}
		controller.service.CreateOrUpdateGauge(params.Name, value)
	case models.Counter:
		value, err := strconv.ParseInt(params.Value, 10, 64)
		if err != nil {
			http.Error(res, "invalid metric value", http.StatusBadRequest)
			return
		}
		controller.service.CreateOrUpdateCounter(params.Name, value)
	default:
		http.Error(res, "invalid metric type", http.StatusBadRequest)
		return
	}


	res.WriteHeader(http.StatusOK)
}
