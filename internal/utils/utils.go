package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/90amper/metmon/internal/logger"
	"github.com/90amper/metmon/internal/models"
)

func PrettyPrint(v interface{}) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return nil
}

func ParseGauge(str string) (res models.Gauge, err error) {
	valf, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}
	return models.Gauge(valf), nil
}

func ParseCounter(str string) (res models.Counter, err error) {
	vali, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return models.Counter(vali), nil
}

// Use: helper function
func Retryer(fn func() error) error {
	// , args ...any
	var err error
	// delay := 0 * time.Second
	for i := 0; i <= 3; i++ {
		if i == 1 {
			logger.Log("attempting to retry $")
		}
		if i > 0 {
			logger.Printf("<%d> ", i)
			delay := time.Duration(i*2-1) * time.Second
			time.Sleep(delay)
		}
		err = fn()
		if err == nil {
			return nil
		}
	}
	return err
}
