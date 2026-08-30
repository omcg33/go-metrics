package service

import (
	"github.com/omcg33/go-metrics/internal/repository"
	"github.com/stretchr/testify/mock"
)

var _ repository.Storage = (*storageMock)(nil)

type storageMock struct {
	mock.Mock
}

func (m *storageMock) UpdateGauge(name string, value float64) {
	m.Called(name, value)
}

func (m *storageMock) UpdateCounter(name string, delta int64) {
	m.Called(name, delta)
}

func (m *storageMock) Gauges() map[string]float64 {
	args := m.Called()
	result, _ := args.Get(0).(map[string]float64)
	return result
}

func (m *storageMock) Counters() map[string]int64 {
	args := m.Called()
	result, _ := args.Get(0).(map[string]int64)
	return result
}
