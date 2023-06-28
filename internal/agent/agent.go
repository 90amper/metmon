package agent

import (
	"time"

	"github.com/90amper/metmon/internal/agent/collector"
	"github.com/90amper/metmon/internal/agent/config"
	"github.com/90amper/metmon/internal/agent/sender"
	"github.com/90amper/metmon/internal/storage"
)

type Agent struct {
	PollInterval   time.Duration
	ReportInterval time.Duration
	Storage        storage.Storager
	Collector      *collector.Collector
	Sender         *sender.Sender
}

func NewAgent() (agent *Agent, err error) {
	var a Agent
	a.PollInterval = time.Duration(config.Config.PollInterval) * time.Second
	a.ReportInterval = time.Duration(config.Config.ReportInterval) * time.Second
	a.Storage = storage.NewStorage(&config.Config)
	a.Collector, err = collector.NewCollector(config.Config, a.Storage)
	if err != nil {
		return nil, err
	}
	a.Sender, err = sender.NewSender(config.Config, a.Storage)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
