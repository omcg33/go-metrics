package main

import (
	"flag"
	"time"

	"github.com/omcg33/go-metrics/internal/agent"
)

var (
	serverAddress *string
	reportInterval *int
	pollInterval *int
)

func init() {
  serverAddress = flag.String("a", "localhost:8080", "Адрес сервера")
  reportInterval = flag.Int("r", 10 , "Частота отправки метрик на сервер")
  pollInterval = flag.Int("p", 2 , "Частота опроса метрик из пакета ")
}


func main() {
	flag.Parse();

	reportTicker := time.NewTicker(time.Duration(*reportInterval) * time.Second)
	pollTicker := time.NewTicker(time.Duration(*pollInterval) * time.Second)
	collector := agent.NewRuntimeCollector()
	service := agent.NewService(*serverAddress)

	defer pollTicker.Stop()
	defer reportTicker.Stop()

	for {
		select {
			case <-pollTicker.C:
				collector.Collect()
			case <-reportTicker.C:
				report := collector.Report()
				service.Report(report)
		}
	}
}