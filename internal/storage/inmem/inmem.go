package inmem

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/90amper/metmon/internal/logger"
	"github.com/90amper/metmon/internal/models"
)

type MemStorage struct {
	models.Store
	cfg *models.Config
}

func NewInMem(cfg *models.Config) *MemStorage {
	return &MemStorage{
		cfg: cfg,
	}
}

func (s *MemStorage) PingDB() error {
	return fmt.Errorf("not supported in memory storage")
}

func (s *MemStorage) BatchAdd(ms []models.Metric) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	var errs []error
	for _, metric := range ms {
		var err error
		switch metric.MType {
		case "gauge":
			err = s.AddGauge(metric.ID, models.Gauge(*metric.Value))
			if err != nil {
				errs = append(errs, err)
			}
		case "counter":
			err = s.AddCounter(metric.ID, models.Counter(*metric.Delta))
			if err != nil {
				errs = append(errs, err)
			}
		default:
			err := fmt.Errorf("unsupported metric type")
			errs = append(errs, err)
			logger.Error(err)
			continue
		}
	}
	return errors.Join(errs...)
}

func (s *MemStorage) AddGauge(name string, value models.Gauge) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
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
	s.Mu.Lock()
	defer s.Mu.Unlock()
	if s.Counters == nil {
		s.Counters = make(models.CounterStore)
	}
	if _, ok := s.Counters[name]; !ok {
		s.Counters[name] = 0
	}
	s.Counters[name] += value

	return nil
}

func (s *MemStorage) TickCounter(name string) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	if s.Counters == nil {
		s.Counters = make(models.CounterStore)
	}
	if _, ok := s.Counters[name]; !ok {
		s.Counters[name] = 0
	}

	s.Counters[name]++

	return nil
}

func (s *MemStorage) Purge() error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Counters = make(models.CounterStore)
	s.Gauges = nil
	return nil
}

func (s *MemStorage) GetAllMetrics() ([]models.Metric, error) {
	var marr []models.Metric

	s.Mu.Lock()
	defer s.Mu.Unlock()
	gauges := s.Gauges
	counters := s.Counters

	for name, values := range gauges {
		for _, value := range values {
			val := float64(value)
			metr := models.Metric{
				ID:    name,
				MType: "gauge",
				Value: &val,
			}
			marr = append(marr, metr)
		}
	}
	for name, value := range counters {
		val := int64(value)
		metr := models.Metric{
			ID:    name,
			MType: "counter",
			Delta: &val,
		}
		marr = append(marr, metr)
	}

	return marr, nil
}

func (s *MemStorage) GetGauges() (models.GaugeStore, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	return s.Gauges, nil
}

func (s *MemStorage) GetCounters() (models.CounterStore, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	return s.Counters, nil
}

func (s *MemStorage) GetCurrentGauge(mName string) (models.Gauge, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	if _, ok := s.Gauges[mName]; !ok {
		return 0, fmt.Errorf("gauge %s not found", mName)
	}
	return s.Gauges[mName][len(s.Gauges[mName])-1], nil
}

func (s *MemStorage) GetCounter(mName string) (models.Counter, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	if _, ok := s.Counters[mName]; !ok {
		return 0, fmt.Errorf("counter %s not found", mName)
	}
	return s.Counters[mName], nil
}

func (s *MemStorage) GetCurrentGauges() (models.GaugeList, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	if s.Gauges == nil {
		return models.GaugeList{}, nil
	}
	list := make(models.GaugeList)
	for mName, mVal := range s.Gauges {
		list[mName] = mVal[len(mVal)-1]
	}
	return list, nil
}

func (s *MemStorage) SaveToFile() error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	json, err := json.Marshal(s)
	if err != nil {
		logger.Log(err.Error())
	}
	err = os.WriteFile(s.cfg.FileStoragePath, json, 0666)
	if err != nil {
		logger.Log(err.Error())
	}

	return nil
}
func (s *MemStorage) LoadFromFile() error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	fmt.Printf("Loading storage from file... ")
	store := &MemStorage{}
	file, err := os.ReadFile(s.cfg.FileStoragePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, store)
	if err != nil {
		return err
	}
	s.Counters = store.Counters
	s.Gauges = store.Gauges
	fmt.Printf("done: %v gauges, %v counters\n", len(store.Gauges), len(store.Counters))
	return nil
}

func (s *MemStorage) Dumper() error {
	fmt.Println("Dumper started")
	for {
		if s.cfg.StoreInterval > 0 {
			time.Sleep(time.Duration(s.cfg.StoreInterval) * time.Second)
		}
		fmt.Printf("Dump storage to file... ")
		err := s.SaveToFile()
		if err != nil {
			return err
		}
		fmt.Printf("done\n")
	}
}
