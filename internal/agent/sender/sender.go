package sender

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/90amper/metmon/internal/logger"
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

func (s *Sender) BatchSend() (err error) {
	basePath := "updates/"

	marr, err := s.storage.GetAllMetrics()
	if err != nil {
		return err
	}

	if len(marr) == 0 {
		return nil
	}

	err = s.Post(basePath, marr)
	if err != nil {
		logger.Error(err)
	}

	err = s.storage.Purge()
	if err != nil {
		logger.Error(err)
	}

	return nil
}

func (s *Sender) Run(wg *sync.WaitGroup) error {
	fmt.Println("Sender started")
	defer wg.Done()
	for {
		err := s.BatchSend()
		if err != nil {
			return err
		}
		time.Sleep(s.reportInterval)
	}
}
