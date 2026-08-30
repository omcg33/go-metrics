package agent

type TService interface {
	Report(Report)
}

type Report struct {
	gauges map[string]float64
	counters map[string]int64
}

type TCollector interface {
	Collect()
	Report() Report
}

type Config struct {
    ServerAddress *string
    ReportInterval *int
	PollInterval *int
}