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

type FileStorageConfig struct {
	StoreInterval time.Duration
	Restore       bool
	FilePath      string
}

type FileStorage struct {
	config  FileStorageConfig
	storage *MemStorage
	logger  *logs.Logger
}

func NewFileStorage(config FileStorageConfig, storage *MemStorage, logger *logs.Logger) (*FileStorage, error) {
	f := &FileStorage{
		config:  config,
		storage: storage,
		logger:  logger,
	}

	if f.config.Restore {
		err := f.restore()

		if err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (f *FileStorage) restore() error {
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

func (f *FileStorage) save() error {
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

func (f *FileStorage) AutoSave() {
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

func (f *FileStorage) UpdateGauge(name string, value float64) error {
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

func (f *FileStorage) UpdateCounter(name string, delta int64) error {
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

func (f *FileStorage) Gauges() map[string]float64 {
	return f.storage.Gauges()
}

func (f *FileStorage) Counters() map[string]int64 {
	return f.storage.Counters()
}
