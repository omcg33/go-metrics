package service

type Storage interface {
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, delta int64)
	Gauges() (map[string]float64 )
	Counters() (map[string]int64)
}