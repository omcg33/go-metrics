package handler

import metrics "github.com/omcg33/go-metrics/internal/service"

type Controller struct {
	service metrics.TMetricsService
}

func NewController(service metrics.TMetricsService) *Controller {
	return &Controller{
		service: service,
	}
}
