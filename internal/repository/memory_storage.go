package repository

var _ TStorage = (*MemStorage)(nil)

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (s *MemStorage) UpdateGauge(name string, value float64) {
	s.gauges[name] = value
}

func (s *MemStorage) UpdateCounter(name string, delta int64) {
	s.counters[name] += delta
}

func (s *MemStorage) Gauges() map[string]float64 {
	return s.gauges
}

func (s *MemStorage) Counters() map[string]int64 {
	return s.counters
}
