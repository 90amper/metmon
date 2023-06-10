package handlers

import (
	"fmt"
	"net/http"

	"github.com/90amper/metmon/internal/storage"
	"github.com/90amper/metmon/internal/utils"
	"github.com/go-chi/chi/v5"
)

type Wrapper struct {
	// config  string
	storage storage.Storager
}

func NewWrapper(storage storage.Storager) (wr *Wrapper, err error) {
	return &Wrapper{
		storage: storage,
	}, nil
}

func (wr *Wrapper) ReceiveMetrics(w http.ResponseWriter, r *http.Request) {
	// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
	mType := chi.URLParam(r, "type")
	mName := chi.URLParam(r, "name")
	mValue := chi.URLParam(r, "value")

	if mName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch mType {
	case "gauge":
		val, err := utils.ParseGauge(mValue)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		wr.storage.AddGauge(mName, val)
	case "counter":
		val, err := utils.ParseCounter(mValue)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		wr.storage.AddCounter(mName, val)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

func (wr *Wrapper) GetCurrentMetric(w http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, "type")
	mName := chi.URLParam(r, "name")
	// logger.Log(mType,mName)
	switch mType {
	case "gauge":
		val, err := wr.storage.GetCurrentGauge(mName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write([]byte(fmt.Sprintf("%v", val)))
	case "counter":
		val, err := wr.storage.GetCounter(mName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write([]byte(fmt.Sprintf("%v", val)))
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

}

func (wr *Wrapper) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("All metrics"))
}
