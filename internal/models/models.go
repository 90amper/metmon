package models

import "sync"

type (
	Gauge   float64
	Counter int64
)

type Store struct {
	Mu       sync.Mutex
	Gauges   GaugeStore   `json:"gauges"`
	Counters CounterStore `json:"counters"`
}

type GaugeStore map[string][]Gauge

type CounterStore map[string]Counter

type GaugeList map[string]Gauge

type Config struct {
	ServerURL       string `env:"ADDRESS"`
	ReportInterval  int    `env:"REPORT_INTERVAL"`
	PollInterval    int    `env:"POLL_INTERVAL"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	PathSeparator   string
	ProjPath        string
	Cleanup         bool
	DatabaseDsn     string `env:"DATABASE_DSN"`
}
