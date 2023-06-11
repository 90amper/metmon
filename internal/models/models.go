package models

type (
	Gauge   float64
	Counter int64
)

type GaugeStore map[string][]Gauge

type CounterStore map[string]Counter

type GaugeList map[string]Gauge

// type CounterList map[string]Counter

type Config struct {
	ServerURL      string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}
