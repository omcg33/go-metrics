package agent

import (
	"log"
	"strings"

	models "github.com/omcg33/go-metrics/internal/model"
	"resty.dev/v3"
)

var _ Service = (*MetricsService)(nil)

type MetricsService struct {
	client *resty.Client
}

func NewService(serverAddress string) *MetricsService {
	if !strings.Contains(serverAddress, "://") {
		serverAddress = "http://" + serverAddress
	}

	return &MetricsService{
		client: resty.New().SetBaseURL(serverAddress),
	}
}

func (service *MetricsService) Close() error {
	return service.client.Close()
}

func (service *MetricsService) Report(report Report) {
	for name, value := range report.gauges {
		log.Printf("Send Gauge metric %s === %f POST to /update", name, value)

		_, err := service.client.R().
			SetBody(models.Metrics{
				ID:    name,
				MType: models.Gauge,
				Value: &value,
			}).
			SetHeader("Content-Type", "application/json").
			Post("/update")

		if err != nil {
			log.Printf("Failed Gauge metric %s === %f POST to /update  with %v", name, value, err)
			panic(err)
		}
	}

	for name, value := range report.counters {

		log.Printf("Send Counter metric %s === %d POST to /update", name, value)

		_, err := service.client.R().
			SetBody(models.Metrics{
				ID:    name,
				MType: models.Counter,
				Delta: &value,
			}).
			SetHeader("Content-Type", "application/json").
			Post("/update")

		if err != nil {
			log.Printf("Failed Counter metric %s === %d POST to /update with %v", name, value, err)
			panic(err)
		}
	}
}
