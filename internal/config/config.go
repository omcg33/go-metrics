package config

import (
	flag "flag"
	log "log"

	env "github.com/caarlos0/env/v6"
)

func NewConfig() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "Адрес сервера")
	flag.Parse()

	err := env.Parse(cfg)

	if err != nil {
		log.Fatal(err)
	}

	return cfg
}
