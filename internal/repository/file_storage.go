package repository

import (
	"encoding/json"
	"os"
	"time"

	"github.com/omcg33/go-metrics/internal/logger"
	models "github.com/omcg33/go-metrics/internal/model"
	"go.uber.org/zap"
)

type FileStorageConfig struct {
	StoreInterval int
	Restore       bool
	FilePath      string
}

type FileStorage struct {
	storage *MemStorage
	config  FileStorageConfig
}

func NewFileStorage(config FileStorageConfig, storage *MemStorage) (*FileStorage, error) {
	f := &FileStorage{
		storage: storage,
		config:  config,
	}

	if f.config.Restore {
		err := f.Restore()

		if err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (f *FileStorage) Restore() error {
	data, err := os.ReadFile(f.config.FilePath)

	if err != nil {
		return err
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

	logger.Log.Info("Metrics restored from file")

	return nil
}

func (f *FileStorage) AutoSave() {
	if f.config.StoreInterval == 0 {
		return
	}

	for range time.Tick(time.Duration(f.config.StoreInterval) * time.Second) {
		err := f.Save()

		if err != nil {
			logger.Log.Fatal("Error saving data to file", zap.Error(err))
		}
	}
}

func (f *FileStorage) Save() error {
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

	logger.Log.Info("Metrics saved to files")

	return os.WriteFile(f.config.FilePath, data, 0666)
}

func (f *FileStorage) UpdateGauge(name string, value float64) {
	f.storage.UpdateGauge(name, value)

	if f.config.StoreInterval == 0 {
		f.Save()
	}
}

func (f *FileStorage) UpdateCounter(name string, delta int64) {
	f.storage.UpdateCounter(name, delta)

	if f.config.StoreInterval == 0 {
		f.Save()
	}
}

func (f *FileStorage) Gauges() map[string]float64 {
	return f.storage.Gauges()
}

func (f *FileStorage) Counters() map[string]int64 {
	return f.storage.Counters()
}
