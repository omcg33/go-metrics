package handler

type MetricsService interface {
	CreateOrUpdateGauge(name string, value float64)
	CreateOrUpdateCounter(name string, value int64)
	Metrics() map[string]interface{}
	Gauge(name string) (float64, bool)
	Counter(name string) (int64, bool)
}