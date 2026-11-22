package repository

import (
	"encoding/json"
	"os"
	"time"

	renameio "github.com/google/renameio"
	logs "github.com/omcg33/go-metrics/internal/logger"
	models "github.com/omcg33/go-metrics/internal/model"
	"go.uber.org/zap"
)

var _ MetricsRepository = (*MetricsFileRepository)(nil)

type Runner interface {
	Run()
}

type MetricsFileRepositoryConfig struct {
	StoreInterval time.Duration
	Restore       bool
	FilePath      string
}

type MetricsFileRepository struct {
	config  MetricsFileRepositoryConfig
	storage *MetricsMemoryRepository
	logger  *logs.Logger
}

func NewMetricsFileRepository(config MetricsFileRepositoryConfig, storage *MetricsMemoryRepository, logger *logs.Logger) (*MetricsFileRepository, error) {
	f := &MetricsFileRepository{
		config:  config,
		storage: storage,
		logger:  logger,
	}

	logger.Info("MetricsFileRepository initialized")

	if f.config.Restore {
		err := f.restore()

		if err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (f *MetricsFileRepository) restore() error {
	data, err := os.ReadFile(f.config.FilePath)

	if err != nil {
		f.logger.Info("Metrics does not restored.", zap.String("filePath", f.config.FilePath))
		return nil
	}

	var metrics []models.Metrics

	if err := json.Unmarshal(data, &metrics); err != nil {
		return err
	}

	for _, metric := range metrics {
		switch metric.MType {
		case models.Gauge:
			f.storage.UpdateGauge(metric.ID, *metric.Value)
		case models.Counter:
			f.storage.UpdateCounter(metric.ID, *metric.Delta)
		}
	}

	f.logger.Info("Metrics restored from file")

	return nil
}

func (f *MetricsFileRepository) save() error {
	gauges := f.storage.Gauges()
	counters := f.storage.Counters()
	metrics := make([]models.Metrics, 0, len(gauges)+len(counters))

	for id, v := range gauges {
		metrics = append(metrics, models.Metrics{ID: id, MType: models.Gauge, Value: &v})
	}
	for id, d := range counters {
		metrics = append(metrics, models.Metrics{ID: id, MType: models.Counter, Delta: &d})
	}

	data, err := json.Marshal(metrics)

	if err != nil {
		return err
	}

	f.logger.Info("Metrics saved to files")

	return renameio.WriteFile(f.config.FilePath, data, 0666)
}

func (f *MetricsFileRepository) AutoSave() {
	if f.config.StoreInterval == 0 {
		return
	}

	for range time.Tick(f.config.StoreInterval) {
		err := f.save()

		if err != nil {
			f.logger.Fatal("Error saving data to file", zap.Error(err))
		}
	}
}

func (f *MetricsFileRepository) Run() { f.AutoSave() }

func (f *MetricsFileRepository) UpdateGauge(name string, value float64) error {
	f.storage.UpdateGauge(name, value)

	if f.config.StoreInterval == 0 {
		err := f.save()

		if err != nil {
			f.logger.Error("Failed to save metrics to file after UpdateGauge()", zap.Error(err))
			return err
		}
	}

	return nil
}

func (f *MetricsFileRepository) UpdateCounter(name string, delta int64) error {
	f.storage.UpdateCounter(name, delta)

	if f.config.StoreInterval == 0 {
		err := f.save()

		if err != nil {
			f.logger.Error("Failed to save metrics to file after UpdateCounter()", zap.Error(err))
			return err
		}
	}

	return nil
}

func (f *MetricsFileRepository) Gauges() map[string]float64 {
	return f.storage.Gauges()
}

func (f *MetricsFileRepository) Counters() map[string]int64 {
	return f.storage.Counters()
}

func (f *MetricsFileRepository) Close() error {
	err := f.save()

	return err
}
