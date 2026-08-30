package handler


type Controller struct {
	service MetricsService
}

func NewController(service MetricsService) *Controller {
	return &Controller{
		service: service,
	}
}
