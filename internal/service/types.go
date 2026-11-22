package service

type MetricsRepository interface {
	UpdateGauge(name string, value float64) error
	UpdateCounter(name string, delta int64) error
	Gauges() map[string]float64
	Counters() map[string]int64
	Close() error
}
