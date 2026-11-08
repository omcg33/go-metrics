package repository

import (
	"maps"
	"sync"
)

type MemStorage struct {
	mu       sync.RWMutex
	gauges   map[string]float64
	counters map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (s *MemStorage) UpdateGauge(name string, value float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.gauges[name] = value
	return nil
}

func (s *MemStorage) UpdateCounter(name string, delta int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counters[name] += delta
	return nil
}

func (s *MemStorage) Gauges() map[string]float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return maps.Clone(s.gauges)
}

func (s *MemStorage) Counters() map[string]int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return maps.Clone(s.counters)
}
