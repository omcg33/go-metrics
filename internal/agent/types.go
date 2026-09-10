package agent

type Service interface {
	Report(Report)
}

type Report struct {
	gauges   map[string]float64
	counters map[string]int64
}

type Collector interface {
	Collect()
	Report() Report
}

type Config struct {
	ServerAddress  string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}
