package inmem

import (
	"github.com/90amper/metmon/internal/models"
)

type MemStorage struct {
	Gauges   models.GaugeStore
	Counters models.CounterStore
}

func NewInMem() *MemStorage {
	return &MemStorage{}
}

func (s *MemStorage) AddGauge(name string, value models.Gauge) error {
	if s.Gauges == nil {
		s.Gauges = make(models.GaugeStore)
	}
	if _, ok := s.Gauges[name]; !ok {
		s.Gauges[name] = []models.Gauge{}
	}
	s.Gauges[name] = append(s.Gauges[name], value)
	return nil
}

func (s *MemStorage) TickCounter(name string) error {
	if s.Counters == nil {
		s.Counters = make(models.CounterStore)
	}
	if _, ok := s.Counters[name]; !ok {
		s.Counters[name] = 0
	} else {
		s.Counters[name]++
	}
	return nil
}

func (s *MemStorage) CleanGauges() error {
	s.Gauges = nil
	return nil
}

func (s *MemStorage) GetGauges() (models.GaugeStore, error) {
	return s.Gauges, nil
}

func (s *MemStorage) GetCounters() (models.CounterStore, error) {
	return s.Counters, nil
}
