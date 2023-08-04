package sender

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/90amper/metmon/internal/logger"
	"github.com/90amper/metmon/internal/models"
	"github.com/90amper/metmon/internal/storage"
	"github.com/90amper/metmon/pkg/hasher"

	"github.com/davecgh/go-spew/spew"
)

type Sender struct {
	client         *http.Client
	destURL        string
	reportInterval time.Duration
	storage        storage.Storager
	hashKey        string
	hashAlg        string
}

func NewSender(config models.Config, storage storage.Storager) (*Sender, error) {
	reportInterval := time.Duration(config.ReportInterval) * time.Second
	return &Sender{
		client:         &http.Client{},
		destURL:        "http://" + config.ServerURL,
		reportInterval: reportInterval,
		storage:        storage,
		hashKey:        config.HashKey,
		hashAlg:        config.HashAlg,
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

	var (
		jbuf, gbuf bytes.Buffer
		bodyHash   []byte
		err        error
	)

	json.NewEncoder(&jbuf).Encode(body)

	gz := gzip.NewWriter(&gbuf)
	gz.Write(jbuf.Bytes())
	gz.Close()

	hbuf, _ := io.ReadAll(&gbuf)
	sendBuf := io.NopCloser(bytes.NewBuffer(hbuf))

	if s.hashKey != "" {
		bodyHash, err = hasher.Hash((hbuf), []byte(s.hashKey), s.hashAlg)
		if err != nil {
			return err
		}
	}

	bodyHashB64 := base64.StdEncoding.EncodeToString(bodyHash)

	req, err := http.NewRequest(http.MethodPost, s.destURL+"/"+path, sendBuf)
	if err != nil {
		return err
	}
	logger.Debug("Body hash: %s", bodyHashB64)
	if s.hashKey != "" {
		req.Header.Set("HashSHA256", string(bodyHashB64))
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
