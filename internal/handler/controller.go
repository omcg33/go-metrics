package handler

import metrics "github.com/omcg33/go-metrics/internal/service"

type Controller struct {
	service metrics.Service
}

func NewController(service metrics.Service) *Controller {
	return &Controller{
		service: service,
	}
}
