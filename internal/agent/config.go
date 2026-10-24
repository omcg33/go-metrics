package agent

import (
	flag "flag"
	log "log"

	env "github.com/caarlos0/env/v6"
)

func NewConfig() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "Адрес сервера")
	flag.IntVar(&cfg.ReportInterval, "r", 10, "Частота отправки метрик на сервер")
	flag.IntVar(&cfg.PollInterval, "p", 2, "Частота опроса метрик из пакета ")
	flag.Parse()

	err := env.Parse(cfg)

	if err != nil {
		log.Fatal(err)
	}

	return cfg
}
