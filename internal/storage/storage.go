package storage

import (
	"github.com/90amper/metmon/internal/models"
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
	SaveToFile() error
	LoadFromFile() error
	Dumper() error
	PingDB() error
}

// func NewStorage(cfg *models.Config) Storager {
// 	// return inmem.NewInMem(cfg)
// 	return sqlbase.NewSqlBase(cfg)
// }
