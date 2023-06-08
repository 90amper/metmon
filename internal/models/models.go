package models

type (
	Gauge   float64
	Counter int64
)

type AgentConfig = struct {
	PollInterval   string
	ReportInterval string
	DestUrl        string
}
