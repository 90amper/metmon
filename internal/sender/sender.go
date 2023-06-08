package sender

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/90amper/metmon/internal/config"
	"github.com/90amper/metmon/internal/storage"
)

type Sender struct {
	client         *http.Client
	destURL        string
	reportInterval time.Duration
	storage        *storage.Storage
}

func NewSender(config config.AgentConfig, storage *storage.Storage) (*Sender, error) {
	reportInterval, err := time.ParseDuration(config.ReportInterval)
	if err != nil {
		return nil, err
	}
	return &Sender{
		client:         &http.Client{},
		destURL:        config.DestURL,
		reportInterval: reportInterval,
		storage:        storage,
	}, nil
}

func (s *Sender) Post(path string) error {
	fmt.Println(path)
	req, err := http.NewRequest(http.MethodPost, s.destURL+"/"+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")

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
	return nil
}

func (s *Sender) SendGauges() error {
	basePath := "update/gauge"

	for name, values := range s.storage.Gauge {
		namePath := basePath + "/" + name
		for _, value := range values {
			path := namePath + "/" + fmt.Sprintf("%f", value)
			err := s.Post(path)
			if err != nil {
				// return err
				fmt.Println(err.Error())
			}
		}
	}
	return nil
}

func (s *Sender) SendCounters() error {
	basePath := "update/counter"
	for name, value := range s.storage.Counter {
		namePath := basePath + "/" + name
		path := namePath + "/" + fmt.Sprintf("%d", value)
		err := s.Post(path)
		if err != nil {
			// return err
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
