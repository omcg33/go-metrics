package agent

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRuntimeCollector_NotNil(t *testing.T) {
	assert.NotNil(t, NewRuntimeCollector())
}

func TestNewRuntimeCollector_GaugesEmpty(t *testing.T) {
	assert.Empty(t, NewRuntimeCollector().gauges)
}

func TestNewRuntimeCollector_CountersEmpty(t *testing.T) {
	assert.Empty(t, NewRuntimeCollector().counters)
}

func TestNewRuntimeCollector_GaugeReaders(t *testing.T) {
	c := NewRuntimeCollector()
	ms := runtime.MemStats{
		Alloc:         1,
		BuckHashSys:   2,
		Frees:         3,
		GCCPUFraction: 0.25,
	}

	tests := []struct {
		name string
		want float64
	}{
		{name: "Alloc", want: 1},
		{name: "BuckHashSys", want: 2},
		{name: "Frees", want: 3},
		{name: "GCCPUFraction", want: 0.25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, c.gaugeReaders[tt.name](ms))
		})
	}
}


func TestRuntimeCollector_CollectIncrementsPollCount(t *testing.T) {
	c := NewRuntimeCollector()
	c.Collect()
	c.Collect()

	assert.Equal(t, int64(2), c.counters["PollCount"])
}

func TestRuntimeCollector_CollectRandomValue(t *testing.T) {
	c := NewRuntimeCollector()
	c.Collect()

	assert.True(t, c.gauges["RandomValue"] >= 0 && c.gauges["RandomValue"] < 1)
}

func TestRuntimeCollector_ReportGauges(t *testing.T) {
	c := NewRuntimeCollector()
	c.Collect()

	assert.Equal(t, c.gauges, c.Report().gauges)
}

func TestRuntimeCollector_ReportCounters(t *testing.T) {
	c := NewRuntimeCollector()
	c.Collect()

	assert.Equal(t, c.counters, c.Report().counters)
}
