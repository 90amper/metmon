package storage

import (
	"github.com/90amper/metmon/internal/models"
	"github.com/90amper/metmon/internal/storage/inmem"
)

type Storager interface {
	AddGauge(name string, value models.Gauge) error
	AddCounter(name string, value models.Counter) error
	TickCounter(name string) error
	CleanGauges() error
	ResetCounters() error
	GetGauges() (models.GaugeStore, error)
	GetCounters() (models.CounterStore, error)
	GetCurrentGauge(mName string) (models.Gauge, error)
	GetCurrentGauges() (models.GaugeList, error)
	GetCounter(mName string) (models.Counter, error)

	// GetCurrentCounterMetric(mName string) (models.Counter, error)
	// GetCurrentMetrics() (models.GaugeStore, models.CounterStore, error)
}

func NewStorage() Storager {
	return inmem.NewInMem()
}
