package storage

import (
	"github.com/90amper/metmon/internal/models"
	"github.com/90amper/metmon/internal/storage/inmem"
)

type Storager interface {
	AddGauge(name string, value models.Gauge) error
	TickCounter(name string) error
	CleanGauges() error
	GetGauges() (models.GaugeStore, error)
	GetCounters() (models.CounterStore, error)
}

func NewStorage() Storager {
	return inmem.NewInMem()
}
