package storage

import (
	"github.com/90amper/metmon/internal/models"
)

type Storager interface {
	AddGauge(name string, value models.Gauge) error
	AddCounter(name string, value models.Counter) error
	BatchAdd([]models.Metric) error
	TickCounter(name string) error
	Purge() error
	GetAllMetrics() ([]models.Metric, error)
	GetGauges() (models.GaugeStore, error)
	GetCounters() (models.CounterStore, error)
	GetCurrentGauge(mName string) (models.Gauge, error)
	GetCurrentGauges() (models.GaugeList, error)
	GetCounter(mName string) (models.Counter, error)
	SaveToFile() error
	LoadFromFile() error
	Dumper() error
	PingDB() error
}
