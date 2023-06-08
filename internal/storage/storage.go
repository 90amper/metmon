package storage

import (
	"github.com/90amper/metmon/internal/models"
)

type GaugeStore map[string][]models.Gauge
type CounterStore map[string]models.Counter

type Storage struct {
	Gauge   GaugeStore
	Counter CounterStore
}

func NewStorage() *Storage {
	storage := Storage{
		Counter: make(CounterStore),
		Gauge:   make(GaugeStore),
	}
	return &storage
}

func (s *Storage) AddGauge(name string, value models.Gauge) error {
	if _, ok := s.Gauge[name]; !ok {
		s.Gauge[name] = []models.Gauge{}
	} else {
		s.Gauge[name] = append(s.Gauge[name], value)
	}
	return nil
}

func (s *Storage) TickCounter(name string) error {
	if _, ok := s.Counter[name]; !ok {
		s.Counter[name] = 0
	} else {
		s.Counter[name]++
	}
	return nil
}
