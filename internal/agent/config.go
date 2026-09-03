package agent

import (
	"flag"
)

func NewConfig() *Config {
	serverAddress := flag.String("a", "localhost:8080", "Адрес сервера")
  	reportInterval := flag.Int("r", 10 , "Частота отправки метрик на сервер")
  	pollInterval := flag.Int("p", 2 , "Частота опроса метрик из пакета ")
	flag.Parse();
	
    return &Config{
        ServerAddress: serverAddress,
		ReportInterval: reportInterval,
		PollInterval: pollInterval,
    }
}