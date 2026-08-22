package handler

import metrics "github.com/omcg33/go-metrics/internal/service"

type Controller struct {
	service metrics.MetricsService
}

func NewController(service metrics.MetricsService) *Controller {
	return &Controller{
		service: service,
	}
}
