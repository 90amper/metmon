package models

type (
	Gauge   float64
	Counter int64
)

type GaugeStore map[string][]Gauge

type CounterStore map[string]Counter

type GaugeList map[string]Gauge

// type CounterList map[string]Counter

type CmdFlags struct {
	ServerUrl      string
	ReportInterval int
	PollInterval   int
}
