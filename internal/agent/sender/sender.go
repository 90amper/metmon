package sender

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/90amper/metmon/internal/models"
	"github.com/90amper/metmon/internal/storage"
	"github.com/davecgh/go-spew/spew"
)

type Sender struct {
	client         *http.Client
	destURL        string
	reportInterval time.Duration
	storage        storage.Storager
}

func NewSender(config models.Config, storage storage.Storager) (*Sender, error) {
	reportInterval := time.Duration(config.ReportInterval) * time.Second
	return &Sender{
		client:         &http.Client{},
		destURL:        "http://" + config.ServerURL,
		reportInterval: reportInterval,
		storage:        storage,
	}, nil
}

func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}
	return b.Bytes(), nil
}

func (s *Sender) Post(path string, body interface{}) error {
	spew.Printf("%#v\n", body)

	var jbuf, gbuf bytes.Buffer
	json.NewEncoder(&jbuf).Encode(body)

	gz := gzip.NewWriter(&gbuf)
	gz.Write(jbuf.Bytes())
	gz.Close()

	req, err := http.NewRequest(http.MethodPost, s.destURL+"/"+path, &gbuf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("post error: %d %s", resp.StatusCode, resp.Status)
	}

	return nil
}

func (s *Sender) SendStore() (err error) {
	err = s.SendGauges()
	if err != nil {
		return
	}
	err = s.storage.CleanGauges()
	if err != nil {
		return
	}

	err = s.SendCounters()
	if err != nil {
		return
	}
	err = s.storage.ResetCounters()
	if err != nil {
		return
	}

	return nil
}

func (s *Sender) SendGauges() error {
	basePath := "update"
	gauges, err := s.storage.GetGauges()
	if err != nil {
		return err
	}
	if gauges == nil {
		return nil
	}
	for name, values := range gauges {
		for _, value := range values {
			val := float64(value)
			metr := models.Metric{
				ID:    name,
				MType: "gauge",
				Value: &val,
			}
			err := s.Post(basePath, metr)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
	return nil
}

func (s *Sender) SendCounters() error {
	basePath := "update"
	counters, err := s.storage.GetCounters()
	if err != nil {
		return err
	}
	if counters == nil {
		return nil
	}
	for name, value := range counters {
		val := int64(value)
		metr := models.Metric{
			ID:    name,
			MType: "counter",
			Delta: &val,
		}
		err := s.Post(basePath, metr)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	return nil
}

func (s *Sender) Run(wg *sync.WaitGroup) error {
	fmt.Println("Sender started")
	defer wg.Done()
	for {
		err := s.SendStore()
		if err != nil {
			return err
		}
		time.Sleep(s.reportInterval)
	}
}
