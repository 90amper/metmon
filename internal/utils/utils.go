package utils

import (
	"encoding/json"
	"fmt"
	"strconv"

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
