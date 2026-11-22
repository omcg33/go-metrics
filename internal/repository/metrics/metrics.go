package repository

import (
	"errors"

	"github.com/omcg33/go-metrics/internal/config"
	logs "github.com/omcg33/go-metrics/internal/logger"
)

type factory struct {
	match func(*config.Config) bool
	build func(*config.Config, *logs.Logger) (MetricsRepository, error)
}

func NewMetricsRepository(cfg *config.Config, logger *logs.Logger) (MetricsRepository, error) {
	factories := []factory{
		{
			match: func(c *config.Config) bool { return c.DatabaseDsn != "" },
			build: buildPostgresRepository,
		},
		{
			match: func(c *config.Config) bool { return c.FileStoragePath != "" },
			build: buildFileRepository,
		},
		{
			match: func(*config.Config) bool { return true },
			build: buildMemoryRepository,
		},
	}

	for _, f := range factories {
		if f.match(cfg) {
			return f.build(cfg, logger)
		}
	}

	return nil, errors.New("no storage factory matched")
}

func buildPostgresRepository(cfg *config.Config, logger *logs.Logger) (MetricsRepository, error) {
	return NewMetricsPostgressRepository(
		MetricsPostgressRepositoryConfig{Dsn: cfg.DatabaseDsn},
		logger,
	)
}

func buildFileRepository(cfg *config.Config, logger *logs.Logger) (MetricsRepository, error) {
	return NewMetricsFileRepository(
		MetricsFileRepositoryConfig{
			FilePath:      cfg.FileStoragePath,
			Restore:       cfg.Restore,
			StoreInterval: cfg.StoreInterval,
		},
		NewMetricsMemoryRepository(),
		logger,
	)
}

func buildMemoryRepository(*config.Config, *logs.Logger) (MetricsRepository, error) {
	return NewMetricsMemoryRepository(), nil
}
