package service

type MetricsService struct {
	repository Storage
}

func NewService(repository Storage) *MetricsService {
	return &MetricsService{repository: repository}
}

func (s *MetricsService) Gauge(name string) (float64, bool) {
	value, ok := s.repository.Gauges()[name]
	return value, ok
}

func (s *MetricsService) Counter(name string) (int64, bool) {
	value, ok := s.repository.Counters()[name]
	return value, ok
}

func (s *MetricsService) CreateOrUpdateGauge(name string, value float64) {
	s.repository.UpdateGauge(name, value)
}

func (s *MetricsService) CreateOrUpdateCounter(name string, value int64) {
	s.repository.UpdateCounter(name, value)
}

func (s *MetricsService) Metrics() map[string]interface{} {
	metrics := make(map[string]interface{})

	gauges := s.repository.Gauges()
	counters := s.repository.Counters()

	metrics["gauges"] = gauges
	metrics["counters"] = counters
	
	
	return metrics
}