package config

import (
	flag "flag"
	log "log"

	"github.com/caarlos0/env/v6"
)

func NewConfig() *Config {
	var cfg Config

	err := env.Parse(&cfg)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(cfg)

	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "Адрес сервера")
	flag.Parse()

	return &cfg
}
