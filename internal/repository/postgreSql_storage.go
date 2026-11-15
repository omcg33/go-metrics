package repository

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	logs "github.com/omcg33/go-metrics/internal/logger"
	"go.uber.org/zap"
)

type PostgreSQLStorageConfig struct {
	Dsn string
}

type PostgreSQLStorage struct {
	config PostgreSQLStorageConfig
	logger *logs.Logger
	db     *sql.DB
}

func NewPostgreSQLStorage(config PostgreSQLStorageConfig, logger *logs.Logger) (*PostgreSQLStorage, error) {
	db, err := sql.Open("pgx", config.Dsn)
	if err != nil {
		logger.Fatal("Failed to open DB", zap.Error(err))
	}

	storage := &PostgreSQLStorage{
		config: config,
		logger: logger,
		db:     db,
	}

	return storage, nil
}

func (s *PostgreSQLStorage) Ping(ctx context.Context) error { return s.db.PingContext(ctx) }
func (s *PostgreSQLStorage) Close() error                   { return s.db.Close() }
