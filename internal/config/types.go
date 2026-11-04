package config

import "time"

type Config struct {
	ServerAddress    string `env:"ADDRESS"`
	StoreInterval    time.Duration
	StoreIntervalInt int    `env:"STORE_INTERVAL"`
	FileStoragePath  string `env:"FILE_STORAGE_PATH"`
	Restore          bool   `env:"RESTORE"`
}
