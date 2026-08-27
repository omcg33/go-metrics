package agent

import (
	"fmt"
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
		fmt.Printf("Send POST to /update/gauge/%s/%f\n", name, value)

		_, err := service.client.R().
			SetPathParams(map[string]string{
				"name": name,
				"value": strconv.FormatFloat(value, 'f', 6, 64),
			}).
			SetHeader("Content-Type", "text/plain").
			Post("/update/gauge/{name}/{value}")

		if(err != nil) {
			fmt.Errorf("Failed POST to /update/gauge/%s/%f with %w", name, value, err)
			panic(err)
		}
	}
	
	for name, value := range report.counters {

		fmt.Printf("Send POST to /update/counter/%s/%d\n", name, value)

		_, err := service.client.R().
			SetPathParams(map[string]string{
				"name": name,
				"value": strconv.FormatInt(value, 10),
			}).
			SetHeader("Content-Type", "text/plain").
			Post("/update/counter/{name}/{value}")

		if(err != nil) {
			fmt.Errorf("Failed POST to /update/counter/%s/%d with %w",  name, value, err)
			panic(err)
		}
	}
}
