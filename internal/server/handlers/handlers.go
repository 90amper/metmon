package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/90amper/metmon/internal/models"
	"github.com/90amper/metmon/internal/storage"
	"github.com/90amper/metmon/internal/utils"
	"github.com/go-chi/chi/v5"
)

type Wrapper struct {
	// config  string
	storage  storage.Storager
	htmlPath string
}

func NewWrapper(storage storage.Storager, htmlPath string) (wr *Wrapper, err error) {
	return &Wrapper{
		storage:  storage,
		htmlPath: htmlPath,
	}, nil
}

func (wr *Wrapper) ReceiveMetrics(w http.ResponseWriter, r *http.Request) {
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
	templFile, err := os.ReadFile(wr.htmlPath + "\\index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	templ, err := template.New("allMetrics").Parse(string(templFile))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	gauges, err := wr.storage.GetCurrentGauges()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	counters, err := wr.storage.GetCounters()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	data := struct {
		Gauges   models.GaugeList
		Counters models.CounterStore
	}{
		Gauges:   gauges,
		Counters: counters,
	}

	err = templ.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

}
