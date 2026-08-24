package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMemStorage_NotNil(t *testing.T) {
	assert.NotNil(t, NewMemStorage())
}

func TestNewMemStorage_GaugesEmpty(t *testing.T) {
	assert.Empty(t, NewMemStorage().Gauges())
}

func TestNewMemStorage_CountersEmpty(t *testing.T) {
	assert.Empty(t, NewMemStorage().Counters())
}

func TestUpdateGauge_SetsValue(t *testing.T) {
	s := NewMemStorage()
	s.UpdateGauge("Alloc", 1.5)

	assert.Equal(t, 1.5, s.Gauges()["Alloc"])
}

func TestUpdateGauge_OverwritesValue(t *testing.T) {
	s := NewMemStorage()
	s.UpdateGauge("Alloc", 1.5)
	s.UpdateGauge("Alloc", 2.5)

	assert.Equal(t, 2.5, s.Gauges()["Alloc"])
}

func TestUpdateCounter_AddsDelta(t *testing.T) {
	s := NewMemStorage()
	s.UpdateCounter("PollCount", 10)

	assert.Equal(t, int64(10), s.Counters()["PollCount"])
}

func TestUpdateCounter_Increments(t *testing.T) {
	s := NewMemStorage()
	s.UpdateCounter("PollCount", 10)
	s.UpdateCounter("PollCount", 5)

	assert.Equal(t, int64(15), s.Counters()["PollCount"])
}

func TestGauges_ReturnsAll(t *testing.T) {
	s := NewMemStorage()
	s.UpdateGauge("Alloc", 1)
	s.UpdateGauge("Sys", 2)

	assert.Equal(t, map[string]float64{"Alloc": 1, "Sys": 2}, s.Gauges())
}

func TestCounters_ReturnsAll(t *testing.T) {
	s := NewMemStorage()
	s.UpdateCounter("PollCount", 1)
	s.UpdateCounter("Custom", 2)

	assert.Equal(t, map[string]int64{"PollCount": 1, "Custom": 2}, s.Counters())
}
