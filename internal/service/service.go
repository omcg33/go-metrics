package service

type Service interface {
	CreateOrUpdateGauge(name string, value float64)
	CreateOrUpdateCounter(name string, value int64)
	Metrics() map[string]interface{}
}