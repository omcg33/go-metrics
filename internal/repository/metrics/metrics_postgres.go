package repository

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	logs "github.com/omcg33/go-metrics/internal/logger"
	"go.uber.org/zap"
)

var _ MetricsRepository = (*MetricsPostgressRepository)(nil)

type MetricsPostgressRepositoryConfig struct {
	Dsn string
}

type MetricsPostgressRepository struct {
	config MetricsPostgressRepositoryConfig
	logger *logs.Logger
	db     *sql.DB
}

func NewMetricsPostgressRepository(config MetricsPostgressRepositoryConfig, logger *logs.Logger) (*MetricsPostgressRepository, error) {
	db, err := sql.Open("pgx", config.Dsn)
	if err != nil {
		logger.Fatal("Failed to open DB", zap.Error(err))
	}

	storage := &MetricsPostgressRepository{
		config: config,
		logger: logger,
		db:     db,
	}

	return storage, nil
}

// shortcut: no SQL yet, replace before relying on DATABASE_DSN
func (s *MetricsPostgressRepository) UpdateGauge(string, float64) error { return nil }
func (s *MetricsPostgressRepository) UpdateCounter(string, int64) error { return nil }
func (s *MetricsPostgressRepository) Gauges() map[string]float64        { return map[string]float64{} }
func (s *MetricsPostgressRepository) Counters() map[string]int64        { return map[string]int64{} }

func (s *MetricsPostgressRepository) Ping(ctx context.Context) error { return s.db.PingContext(ctx) }
func (s *MetricsPostgressRepository) Close() error                   { return s.db.Close() }
