package inmem

import (
	"fmt"

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

func (s *MemStorage) AddCounter(name string, value models.Counter) error {
	if s.Counters == nil {
		s.Counters = make(models.CounterStore)
	}
	if _, ok := s.Counters[name]; !ok {
		s.Counters[name] = 0
	} else {
		s.Counters[name] += value
	}
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

func (s *MemStorage) GetCurrentGauge(mName string) (models.Gauge, error) {
	if _, ok := s.Gauges[mName]; !ok {
		return 0, fmt.Errorf("gauge %s not found", mName)
	}
	return s.Gauges[mName][len(s.Gauges[mName])-1], nil
}

func (s *MemStorage) GetCounter(mName string) (models.Counter, error) {
	if _, ok := s.Counters[mName]; !ok {
		return 0, fmt.Errorf("counter %s not found", mName)
	}
	return s.Counters[mName], nil
}

func (s *MemStorage) GetCurrentGauges() (models.GaugeList, error) {
	if s.Gauges == nil {
		return models.GaugeList{}, nil
	}
	list := make(models.GaugeList)
	for mName, mVal := range s.Gauges {
		list[mName] = mVal[len(mVal)-1]
	}
	return list, nil
}
