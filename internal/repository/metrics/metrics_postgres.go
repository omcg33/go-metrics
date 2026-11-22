package repository

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	logs "github.com/omcg33/go-metrics/internal/logger"
	"github.com/omcg33/go-metrics/internal/migrations"
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
	if err := migrations.Up(config.Dsn); err != nil {
		return nil, err
	}

	db, err := sql.Open("pgx", config.Dsn)
	if err != nil {
		return nil, err
	}

	storage := &MetricsPostgressRepository{
		config: config,
		logger: logger,
		db:     db,
	}

	logger.Info("MetricsPostgressRepository initialized")
	return storage, nil
}

// shortcut: rows are not stored yet, replace before relying on DATABASE_DSN for metrics
func (s *MetricsPostgressRepository) UpdateGauge(string, float64) error { return nil }
func (s *MetricsPostgressRepository) UpdateCounter(string, int64) error { return nil }
func (s *MetricsPostgressRepository) Gauges() map[string]float64        { return map[string]float64{} }
func (s *MetricsPostgressRepository) Counters() map[string]int64        { return map[string]int64{} }

func (s *MetricsPostgressRepository) Ping(ctx context.Context) error { return s.db.PingContext(ctx) }
func (s *MetricsPostgressRepository) Close() error                   { return s.db.Close() }
