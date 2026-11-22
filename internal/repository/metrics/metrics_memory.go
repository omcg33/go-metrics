package repository

import (
	"maps"
	"sync"

	logs "github.com/omcg33/go-metrics/internal/logger"
)

var _ MetricsRepository = (*MetricsMemoryRepository)(nil)

type MetricsMemoryRepository struct {
	logger   *logs.Logger
	mu       sync.RWMutex
	gauges   map[string]float64
	counters map[string]int64
}

func NewMetricsMemoryRepository(logger *logs.Logger) *MetricsMemoryRepository {
	logger.Info("MetricsMemoryRepository initialized")

	return &MetricsMemoryRepository{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
		logger:   logger,
	}
}

func (s *MetricsMemoryRepository) UpdateGauge(name string, value float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.gauges[name] = value
	return nil
}

func (s *MetricsMemoryRepository) UpdateCounter(name string, delta int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counters[name] += delta
	return nil
}

func (s *MetricsMemoryRepository) Gauges() map[string]float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return maps.Clone(s.gauges)
}

func (s *MetricsMemoryRepository) Counters() map[string]int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return maps.Clone(s.counters)
}

func (s *MetricsMemoryRepository) Close() error {
	return nil
}
