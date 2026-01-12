package config

import "sync/atomic"

type Config struct {
	Metrics *Metrics
}

type Metrics struct {
	fileserverHits atomic.Int32
}

func NewConfig() *Config {
	return &Config{
		Metrics: &Metrics{},
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


