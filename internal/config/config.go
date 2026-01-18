package config

import (
	"os"
	"sync/atomic"
)

type Config struct {
	Metrics *Metrics
	Database *Database
	JWTSecret string
	PolkaAPIKey string
}

type Database struct {
	DbURL string
}

type Metrics struct {
	fileserverHits atomic.Int32
}

func NewConfig() *Config {
	return &Config{
		Metrics: &Metrics{},
		Database: &Database{
			DbURL: os.Getenv("DB_URL"),
		},
		JWTSecret: os.Getenv("JWT_SECRET"),
		PolkaAPIKey: os.Getenv("POLKA_KEY"),
	}
}

func (m *Metrics) Inc() {
	m.fileserverHits.Add(1)
}

func (m *Metrics) Reset() {
	m.fileserverHits.Store(0)
}

func (m *Metrics) Hits() int32 {
	return m.fileserverHits.Load()
}


