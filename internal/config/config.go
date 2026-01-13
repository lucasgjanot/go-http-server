package config

import (
	"os"
	"sync/atomic"
)

type Config struct {
	Metrics *Metrics
	Database *Database
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


