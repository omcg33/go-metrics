package agent

import (
	"log"
	"strconv"
	"strings"

	"resty.dev/v3"
)

var _ Service = (*MetricsService)(nil)

type MetricsService struct {
	client *resty.Client;
}

func NewService(serverAddress string) *MetricsService {
	if !strings.Contains(serverAddress, "://") {
    	serverAddress = "http://" + serverAddress
}

	return &MetricsService{
		client: resty.New().SetBaseURL(serverAddress),
	}
}

func (service *MetricsService) Report(report Report) {
	for name, value := range report.gauges {
		log.Printf("Send POST to /update/gauge/%s/%f", name, value)

		_, err := service.client.R().
			SetPathParams(map[string]string{
				"name": name,
				"value": strconv.FormatFloat(value, 'f', 6, 64),
			}).
			SetHeader("Content-Type", "text/plain").
			Post("/update/gauge/{name}/{value}")

		if(err != nil) {
			log.Printf("Failed POST to /update/gauge/%s/%f with %v", name, value, err)
			panic(err)
		}
	}
	
	for name, value := range report.counters {

		log.Printf("Send POST to /update/counter/%s/%d", name, value)

		_, err := service.client.R().
			SetPathParams(map[string]string{
				"name": name,
				"value": strconv.FormatInt(value, 10),
			}).
			SetHeader("Content-Type", "text/plain").
			Post("/update/counter/{name}/{value}")

		if(err != nil) {
			log.Printf("Failed POST to /update/counter/%s/%d with %v", name, value, err)
			panic(err)
		}
	}
}
