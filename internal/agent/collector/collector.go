package collector

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/90amper/metmon/internal/models"
	"github.com/90amper/metmon/internal/storage"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type Collector struct {
	Storage      storage.Storager
	PollInterval time.Duration
}

func NewCollector(config models.Config, storage storage.Storager) (*Collector, error) {
	pollInterval := time.Duration(config.PollInterval) * time.Second
	return &Collector{
		Storage:      storage,
		PollInterval: pollInterval,
	}, nil
}

func (c *Collector) Collect() error {
	c.ReadRuntimeMetrics()
	c.ReadAddonMetrics()
	c.ReadComputeMetrics()
	return nil
}

func (c *Collector) ReadRuntimeMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	values := reflect.ValueOf(m)
	numfield := values.NumField()
	for i := 0; i < numfield; i++ {
		fName := reflect.TypeOf(m).Field(i).Name
		fType := reflect.ValueOf(m).Field(i).Type().String()
		fValue := reflect.ValueOf(m).Field(i)
		var value models.Gauge
		if fType == "uint64" || fType == "uint32" {
			value = models.Gauge(float64(fValue.Uint()))
		} else if fType == "float64" {
			value = models.Gauge(fValue.Float())
		} else {
			continue
		}
		c.Storage.AddGauge(fName, value)
	}
}

func (c *Collector) ReadAddonMetrics() {
	c.Storage.TickCounter("PollCount")
	c.Storage.AddGauge("RandomValue", models.Gauge(rand.Float64()))

}

func (c *Collector) ReadComputeMetrics() {
	memStat, _ := mem.VirtualMemory()
	totalMemory := memStat.Total
	freeMemory := memStat.Free

	c.Storage.AddGauge("TotalMemory", models.Gauge(totalMemory))
	c.Storage.AddGauge("FreeMemory", models.Gauge(freeMemory))

	CPUutilization, _ := cpu.Percent(0, true)
	for core, cpu := range CPUutilization {
		c.Storage.AddGauge(fmt.Sprintf("CPUutilization%d", core), models.Gauge(cpu))
	}
}

func (c *Collector) Run(wg *sync.WaitGroup) error {
	fmt.Println("Collector started")
	defer wg.Done()
	for {
		err := c.Collect()
		if err != nil {
			return err
		}
		time.Sleep(c.PollInterval)
	}
}
