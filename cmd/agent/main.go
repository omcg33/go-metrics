package main

import (
	"fmt"
	"time"

	"github.com/omcg33/go-metrics/internal/agent"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {
	pollTicker := time.NewTicker(pollInterval)
	reportTicker := time.NewTicker(reportInterval)
	collector := agent.NewRuntimeCollector()
	defer pollTicker.Stop()
	defer reportTicker.Stop()


	for {
	select {
	case <-pollTicker.C:
		collector.Collect()
	case <-reportTicker.C:
		report := collector.Report()
		fmt.Println(report)
	}
}

	
	
	
	
}