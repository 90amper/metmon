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

// Compress сжимает слайс байт.
func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	// создаём переменную w — в неё будут записываться входящие данные,
	// которые будут сжиматься и сохраняться в bytes.Buffer
	w := gzip.NewWriter(&b)
	// запись данных
	_, err := w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}
	// обязательно нужно вызвать метод Close() — в противном случае часть данных
	// может не записаться в буфер b; если нужно выгрузить все упакованные данные
	// в какой-то момент сжатия, используйте метод Flush()
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}
	// переменная b содержит сжатые данные
	return b.Bytes(), nil
}

// func Compressor(req *http.Request) {
//     var b []byte
// 	req.Body.Read(b)
//     var buf bytes.Buffer
//     g := gzip.NewWriter(&buf)

//     _, err := io.Copy(g, &b)
//     if err != nil {
//         logger.Log(err.Error())
//         return
//     }
// }

func (s *Sender) Post(path string, body interface{}) error {
	// fmt.Println(path)
	spew.Printf("%#v\n", body)
	// json, err := json.Marshal(body)
	// if err != nil {
	// 	logger.Log(err.Error())
	// }
	// gzjson, err := Compress(json)
	// if err != nil {
	// 	logger.Log(err.Error())
	// }

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
		// namePath := basePath + "/" + name
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
