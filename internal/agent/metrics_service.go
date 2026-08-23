package agent

import (
	"fmt"
	"net/http"
)

var _ Service = (*MetricsService)(nil)

type MetricsService struct {
	baseUrl string;
	client *http.Client;
}

func NewService() *MetricsService {
	return &MetricsService{
		baseUrl: "http://localhost:8080",
		client: &http.Client{},
	}
}

func (service *MetricsService) Report(report Report) {
	for name, value := range report.gauges {
		fmt.Printf("Send POST to %s/update/gauge/%s/%f\n", service.baseUrl, name, value)
		_, err := service.client.Post(fmt.Sprintf("%s/update/gauge/%s/%f", service.baseUrl, name, value), "text/plain", nil)

		if(err != nil) {
			fmt.Errorf("Failed POST to %s/update/gauge/%s/%f with %w", service.baseUrl, name, value, err)
			panic(err)
		}
	}
	
	for name, value := range report.counters {
		fmt.Printf("Send POST to %s/update/counter/%s/%d\n", service.baseUrl, name, value)
		_, err := service.client.Post(fmt.Sprintf("%s/update/counter/%s/%d", service.baseUrl, name, value), "text/plain", nil)

		if(err != nil) {
			fmt.Errorf("Failed POST to %s/update/counter/%s/%d with %w", service.baseUrl, name, value, err)
			panic(err)
		}
	}
}
