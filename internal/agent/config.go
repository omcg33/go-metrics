package agent

import (
	flag "flag"

	env "github.com/caarlos0/env/v11"
	logs "github.com/omcg33/go-metrics/internal/logger"
	"go.uber.org/zap"
)

func NewConfig(logger *logs.Logger) *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "Адрес сервера")
	flag.IntVar(&cfg.ReportInterval, "r", 10, "Частота отправки метрик на сервер")
	flag.IntVar(&cfg.PollInterval, "p", 2, "Частота опроса метрик из пакета ")
	flag.Parse()

	err := env.Parse(cfg)

	if err != nil {
		logger.Fatal("Failed to parse envs", zap.Error(err))
	}

	return cfg
}
