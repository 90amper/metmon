package sender

import (
	"fmt"
	"net/http"

	"github.com/90amper/metmon/internal/models"
	"github.com/90amper/metmon/internal/storage"
)

type Sender struct {
	client  *http.Client
	destUrl string
}

func NewSender(config models.AgentConfig) *Sender {
	return &Sender{
		client:  &http.Client{},
		destUrl: config.DestUrl,
	}
}

func (s *Sender) Post(path string) error {
	fmt.Println(path)
	req, err := http.NewRequest(http.MethodPost, s.destUrl+"/"+path, nil)
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

func (s *Sender) SendStore(store storage.Storage) (err error) {
	err = s.SendGauges(store.Gauge)
	if err != nil {
		return
	}
	err = s.SendCounters(store.Counter)
	if err != nil {
		return
	}
	return nil
}

func (s *Sender) SendGauges(gauges storage.GaugeStore) error {
	basePath := "update/gauge"
	for name, values := range gauges {
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

func (s *Sender) SendCounters(counters storage.CounterStore) error {
	basePath := "update/counter"
	for name, value := range counters {
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
