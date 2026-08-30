package handler

import (
	"github.com/stretchr/testify/mock"
)

type serviceMock struct {
	mock.Mock
}

func (m *serviceMock) CreateOrUpdateGauge(name string, value float64) {
	m.Called(name, value)
}

func (m *serviceMock) CreateOrUpdateCounter(name string, value int64) {
	m.Called(name, value)
}

func (m *serviceMock) Gauge(name string) (float64, bool) {
	args := m.Called(name)
	result, _ := args.Get(0).(float64)
	return result, args.Bool(1)
}

func (m *serviceMock) Counter(name string) (int64, bool) {
	args := m.Called(name)
	result, _ := args.Get(0).(int64)
	return result, args.Bool(1)
}

func (m *serviceMock) Metrics() map[string]interface{} {
	args := m.Called()
	result, _ := args.Get(0).(map[string]interface{})
	return result
}
