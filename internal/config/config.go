package config

import (
	"flag"
)

func NewConfig() *Config {
	serverAddress := flag.String("a", "localhost:8080", "Адрес сервера")
	flag.Parse();
	
    return &Config{
        ServerAddress: serverAddress,
    }
}