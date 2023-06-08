package main

import (
	"time"

	"github.com/90amper/metmon/internal/collector"
	"github.com/90amper/metmon/internal/storage"
	"github.com/90amper/metmon/internal/utils"
)

type Agent struct {
	PollInterval   time.Duration
	ReportInterval time.Duration
	Storage        *storage.Storage
	Collector      *collector.Collector
}

type AgentConfig = struct {
	PollInterval   string
	ReportInterval string
}

var config = AgentConfig{
	PollInterval:   "2s",
	ReportInterval: "10s",
}

func NewAgent(config AgentConfig) *Agent {
	var a Agent
	a.PollInterval, _ = time.ParseDuration(config.PollInterval)
	a.ReportInterval, _ = time.ParseDuration(config.ReportInterval)
	a.Storage = storage.NewStorage()
	a.Collector = collector.NewCollector(a.Storage)
	return &a
}

func main() {
	agent := NewAgent(config)
	agent.Collector.Collect()
	agent.Collector.Collect()
	agent.Collector.Collect()

	utils.PrettyPrint(agent.Storage)
}
