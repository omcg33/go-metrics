package handler

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	models "github.com/omcg33/go-metrics/internal/model"
)

type Controller struct {
	service MetricsService
	validator *validator.Validate
}

func NewController(service MetricsService) *Controller {
	validator := validator.New()
	validator.RegisterAlias("metric_type", fmt.Sprintf("oneof=%s %s", models.Gauge, models.Counter))

	return &Controller{
		service: service,
		validator: validator,
	}
}
