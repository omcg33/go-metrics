package agent

type Report struct {
	gauges map[string]float64
	counters map[string]int64
}

type Collector interface {
	Collect()
	Report() Report
}
