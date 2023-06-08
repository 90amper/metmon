package collector

import (
	"math/rand"
	"reflect"
	"runtime"

	"github.com/90amper/metmon/internal/models"
	"github.com/90amper/metmon/internal/storage"
)

type Collector struct {
	Storage *storage.Storage
}

func NewCollector(storage *storage.Storage) *Collector {
	return &Collector{Storage: storage}
}

func (c *Collector) Collect() error {
	c.ReadRuntimeMetrics()
	c.ReadAddonMetrics()
	return nil
}

func (c *Collector) ReadRuntimeMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	values := reflect.ValueOf(m)
	numfield := values.NumField()
	// res := make(map[string]float64)
	for i := 0; i < numfield; i++ {
		fName := reflect.TypeOf(m).Field(i).Name
		fType := reflect.ValueOf(m).Field(i).Type().String()
		fValue := reflect.ValueOf(m).Field(i)
		var value models.Gauge
		if fType == "uint64" || fType == "uint32" {
			// fmt.Printf("%20v\t%-10v\t%v\n", fName, fType, fValue.Uint())
			// res[fName] = float64(fValue.Uint())
			value = models.Gauge(float64(fValue.Uint()))
		} else if fType == "float64" {
			// fmt.Printf("%20v\t%-10v\t%v\n", fName, fType, fValue.Float())
			value = models.Gauge(fValue.Float())
		} else {
			// fmt.Printf("skip %100v\t%v\n", fName, fType)
			continue
		}
		c.Storage.AddGauge(fName, value)
	}
}

func (c *Collector) ReadAddonMetrics() {
	c.Storage.TickCounter("PollCount")
	c.Storage.AddGauge("RandomValue", models.Gauge(rand.Float64()))

}
