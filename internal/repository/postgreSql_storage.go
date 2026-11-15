package repository

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	logs "github.com/omcg33/go-metrics/internal/logger"
	"go.uber.org/zap"
)

type PostgreSqlStorageConfig struct {
	Dsn string
}

type PostgreSqlStorage struct {
	config PostgreSqlStorageConfig
	logger *logs.Logger
	db     *sql.DB
}

func NewPostgreSqlStorage(config PostgreSqlStorageConfig, logger *logs.Logger) (*PostgreSqlStorage, error) {
	db, err := sql.Open("pgx", config.Dsn)
	if err != nil {
		logger.Fatal("Failed to open DB", zap.Error(err))
	}

	storage := &PostgreSqlStorage{
		config: config,
		logger: logger,
		db:     db,
	}

	return storage, nil
}

func (s *PostgreSqlStorage) Ping(ctx context.Context) error { return s.db.PingContext(ctx) }
func (s *PostgreSqlStorage) Close() error                   { return s.db.Close() }
