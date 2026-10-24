package agent

import (
	"log"
	"math/rand"
	"runtime"
);

var _ Collector = (*RuntimeCollector)(nil)

type GaugeReader func(ms runtime.MemStats) float64
type CounterReader func() int64

type RuntimeCollector struct {
	gaugeReaders map[string]GaugeReader
	counterReaders map[string]CounterReader
	gauges   map[string]float64
	counters map[string]int64
}

func NewRuntimeCollector() *RuntimeCollector {
	var gaugeReaders = map[string]GaugeReader{
		"Alloc": func(ms runtime.MemStats) float64 {
			return float64(ms.Alloc)
		},
		"BuckHashSys": func(ms runtime.MemStats) float64 {
			return float64(ms.BuckHashSys)
		},
		"Frees": func(ms runtime.MemStats) float64 {
			return float64(ms.Frees)
		},
		"GCCPUFraction": func(ms runtime.MemStats) float64 {
			return float64(ms.GCCPUFraction)
		},
		"GCSys": func(ms runtime.MemStats) float64 {
			return float64(ms.GCSys)
		},
		"HeapAlloc": func(ms runtime.MemStats) float64 {
			return float64(ms.HeapAlloc)
		},
		"HeapIdle": func(ms runtime.MemStats) float64 {
			return float64(ms.HeapIdle)
		},
		"HeapInuse": func(ms runtime.MemStats) float64 {
			return float64(ms.HeapInuse)
		},
		"HeapObjects": func(ms runtime.MemStats) float64 {
			return float64(ms.HeapObjects)
		},
		"HeapReleased": func(ms runtime.MemStats) float64 {
			return float64(ms.HeapReleased)
		},
		"HeapSys": func(ms runtime.MemStats) float64 {
			return float64(ms.HeapSys)
		},
		"LastGC": func(ms runtime.MemStats) float64 {
			return float64(ms.LastGC)
		},
		"Lookups": func(ms runtime.MemStats) float64 {
			return float64(ms.Lookups)
		},
		"MCacheInuse": func(ms runtime.MemStats) float64 {
			return float64(ms.MCacheInuse)
		},
		"MCacheSys": func(ms runtime.MemStats) float64 {
			return float64(ms.MCacheSys)
		},
		"MSpanInuse": func(ms runtime.MemStats) float64 {
			return float64(ms.MSpanInuse)
		},
		"MSpanSys": func(ms runtime.MemStats) float64 {
			return float64(ms.MSpanSys)
		},
		"Mallocs": func(ms runtime.MemStats) float64 {
			return float64(ms.Mallocs)
		},
		"NextGC": func(ms runtime.MemStats) float64 {
			return float64(ms.NextGC)
		},
		"NumForcedGC": func(ms runtime.MemStats) float64 {
			return float64(ms.NumForcedGC)
		},
		"NumGC": func(ms runtime.MemStats) float64 {
			return float64(ms.NumGC)
		},
		"OtherSys": func(ms runtime.MemStats) float64 {
			return float64(ms.OtherSys)
		},
		"PauseTotalNs": func(ms runtime.MemStats) float64 {
			return float64(ms.PauseTotalNs)
		},
		"StackInuse": func(ms runtime.MemStats) float64 {
			return float64(ms.StackInuse)
		},
		"StackSys": func(ms runtime.MemStats) float64 {
			return float64(ms.StackSys)
		},
		"Sys": func(ms runtime.MemStats) float64 {
			return float64(ms.Sys)
		},
		"TotalAlloc": func(ms runtime.MemStats) float64 {
			return float64(ms.TotalAlloc)
		},
		"RandomValue":func(ms runtime.MemStats) float64 {
			return rand.Float64()
		},
	}

	c := &RuntimeCollector{
		gaugeReaders: gaugeReaders,
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}

	c.counterReaders = map[string]CounterReader{
		"PollCount": func() int64 {
			c.counters["PollCount"]++
			return c.counters["PollCount"]
		},
	}

	return c
}

func (c *RuntimeCollector) Collect() {
	var m runtime.MemStats

	runtime.ReadMemStats(&m)
	
	log.Println("Collecting metrics")

	for name, reader := range c.gaugeReaders {
		c.gauges[name] = reader(m)
	}

	for name, reader := range c.counterReaders {
		c.counters[name] = reader()
	}
}

func (c *RuntimeCollector) Report() Report {
	log.Println("Report created")
	return Report{
		gauges: c.gauges,
		counters: c.counters,
	}
}