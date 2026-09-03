package main

import (
	"time"

	"github.com/omcg33/go-metrics/internal/agent"
)

func main() {
	config := agent.NewConfig();

	reportTicker := time.NewTicker(time.Duration(*config.ReportInterval) * time.Second)
	pollTicker := time.NewTicker(time.Duration(*config.PollInterval) * time.Second)
	collector := agent.NewRuntimeCollector()
	service := agent.NewService(*config.ServerAddress)

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