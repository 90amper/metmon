package main

import (
	"fmt"
	"time"

	"github.com/90amper/metmon/internal/collector"
	"github.com/90amper/metmon/internal/models"
	"github.com/90amper/metmon/internal/sender"
	"github.com/90amper/metmon/internal/storage"
)

type Agent struct {
	PollInterval   time.Duration
	ReportInterval time.Duration
	Storage        *storage.Storage
	Collector      *collector.Collector
	Sender         *sender.Sender
}

var config = models.AgentConfig{
	PollInterval:   "2s",
	ReportInterval: "10s",
	DestUrl:        "http://localhost:8080",
}

func NewAgent(config models.AgentConfig) *Agent {
	var a Agent
	a.PollInterval, _ = time.ParseDuration(config.PollInterval)
	a.ReportInterval, _ = time.ParseDuration(config.ReportInterval)
	a.Storage = storage.NewStorage()
	a.Collector = collector.NewCollector(a.Storage)
	a.Sender = sender.NewSender(config)
	return &a
}

func main() {
	agent := NewAgent(config)
	agent.Collector.Collect()
	agent.Collector.Collect()
	agent.Collector.Collect()

	// utils.PrettyPrint(agent.Storage)
	err := agent.Sender.SendStore(*agent.Storage)
	if err != nil {
		fmt.Println(err.Error())
	}
}
