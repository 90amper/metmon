package agent

import (
	"time"

	"github.com/90amper/metmon/internal/collector"
	"github.com/90amper/metmon/internal/config"
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

func NewAgent(config config.AgentConfig) (agent *Agent, err error) {
	var a Agent
	a.PollInterval, _ = time.ParseDuration(config.PollInterval)
	a.ReportInterval, _ = time.ParseDuration(config.ReportInterval)
	a.Storage = storage.NewStorage()
	a.Collector, err = collector.NewCollector(config, a.Storage)
	if err != nil {
		return nil, err
	}
	a.Sender, err = sender.NewSender(config, a.Storage)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
