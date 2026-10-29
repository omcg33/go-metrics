package main

import (
	"log"
	"time"

	"github.com/omcg33/go-metrics/internal/agent"
	logs "github.com/omcg33/go-metrics/internal/logger"
)

func main() {
	logger, err := logs.NewLogger("Info")

	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	config := agent.NewConfig(logger)

	reportTicker := time.NewTicker(time.Duration(config.ReportInterval) * time.Second)
	pollTicker := time.NewTicker(time.Duration(config.PollInterval) * time.Second)
	collector := agent.NewRuntimeCollector(logger)
	service := agent.NewService(config.ServerAddress, logger)

	defer pollTicker.Stop()
	defer reportTicker.Stop()
	defer service.Close()

	collector.Collect()

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
