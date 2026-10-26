package config

import (
	"errors"
	flag "flag"
	"os"

	env "github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/omcg33/go-metrics/internal/logger"
	"go.uber.org/zap"
)

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		logger.Log.Fatal("Error loading .env file - ", zap.Error(err))
	}

	cfg := &Config{}

	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "Адрес сервера")
	flag.IntVar(&cfg.StoreInterval, "i", 300, "Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	flag.StringVar(&cfg.FileStoragePath, "f", "./temp.txt", "Путь до файла, куда сохраняются текущие значения")
	flag.BoolVar(&cfg.Restore, "r", false, "Следует ли загружать ранее сохранённые значения из указанного файла при старте сервера?")
	flag.Parse()

	err := env.Parse(cfg)

	if err != nil {
		logger.Log.Fatal("Failed to parse config", zap.Error(err))
	}

	return cfg
}
