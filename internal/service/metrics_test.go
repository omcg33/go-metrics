package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewService_NotNil(t *testing.T) {
	assert.NotNil(t, NewService(&storageMock{}))
}

func TestNewService_SetsRepository(t *testing.T) {
	repo := &storageMock{}

	assert.Equal(t, repo, NewService(repo).repository)
}

func TestCreateOrUpdateGauge(t *testing.T) {
	repo := &storageMock{}
	repo.On("UpdateGauge", "Alloc", 1.5).Return()

	NewService(repo).CreateOrUpdateGauge("Alloc", 1.5)

	repo.AssertCalled(t, "UpdateGauge", "Alloc", 1.5)
}

func TestCreateOrUpdateCounter(t *testing.T) {
	repo := &storageMock{}
	repo.On("UpdateCounter", "PollCount", int64(10)).Return()

	NewService(repo).CreateOrUpdateCounter("PollCount", 10)

	repo.AssertCalled(t, "UpdateCounter", "PollCount", int64(10))
}

func TestMetrics_Gauges(t *testing.T) {
	repo := &storageMock{}
	gauges := map[string]float64{"Alloc": 1.5}
	repo.On("Gauges").Return(gauges)
	repo.On("Counters").Return(map[string]int64{})

	got := NewService(repo).Metrics()

	assert.Equal(t, gauges, got["gauges"])
}

func TestMetrics_Counters(t *testing.T) {
	repo := &storageMock{}
	counters := map[string]int64{"PollCount": 10}
	repo.On("Gauges").Return(map[string]float64{})
	repo.On("Counters").Return(counters)

	got := NewService(repo).Metrics()

	assert.Equal(t, counters, got["counters"])
}

func TestGauge(t *testing.T) {
	repo := &storageMock{}
	repo.On("Gauges").Return(map[string]float64{"Alloc": 1.5})

	got, ok := NewService(repo).Gauge("Alloc")

	assert.Equal(t, []any{1.5, true}, []any{got, ok})
}

func TestGauge_NotFound(t *testing.T) {
	repo := &storageMock{}
	repo.On("Gauges").Return(map[string]float64{})

	_, ok := NewService(repo).Gauge("Alloc")

	assert.False(t, ok)
}

func TestCounter(t *testing.T) {
	repo := &storageMock{}
	repo.On("Counters").Return(map[string]int64{"PollCount": 10})

	got, ok := NewService(repo).Counter("PollCount")

	assert.Equal(t, []any{int64(10), true}, []any{got, ok})
}

func TestCounter_NotFound(t *testing.T) {
	repo := &storageMock{}
	repo.On("Counters").Return(map[string]int64{})

	_, ok := NewService(repo).Counter("PollCount")

	assert.False(t, ok)
}
