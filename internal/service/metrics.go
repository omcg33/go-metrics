package service

import (
	"github.com/omcg33/go-metrics/internal/repository"
)

var _ MetricsService = (*Service)(nil)

type Service struct {
	repository repository.Storage
}

func NewService(repository repository.Storage) *Service {
	return &Service{repository: repository}
}

func (s *Service) CreateOrUpdateGauge(name string, value float64) {
	s.repository.UpdateGauge(name, value)
}

func (s *Service) CreateOrUpdateCounter(name string, value int64) {
	s.repository.UpdateCounter(name, value)
}

func (s *Service) Metrics() map[string]interface{} {
	metrics := make(map[string]interface{})

	gauges := s.repository.Gauges()
	counters := s.repository.Counters()

	metrics["gauges"] = gauges
	metrics["counters"] = counters
	
	
	return metrics
}