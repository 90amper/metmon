package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/90amper/metmon/internal/logger"
	"github.com/90amper/metmon/internal/models"
	"github.com/90amper/metmon/internal/storage"
	"github.com/90amper/metmon/internal/utils"
	"github.com/go-chi/chi/v5"
)

type MMHandler struct {
	storage storage.Storager
	fsPath  string
}

func NewHandler(storage storage.Storager, fsPath string) (hl *MMHandler, err error) {
	return &MMHandler{
		storage: storage,
		fsPath:  fsPath,
	}, nil
}

func (hl *MMHandler) ReceiveMetrics(w http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, "type")
	mName := chi.URLParam(r, "name")
	mValue := chi.URLParam(r, "value")

	if mName == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Metric name not specified!"))
		return
	}

	switch mType {
	case "gauge":
		val, err := utils.ParseGauge(mValue)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad gauge value: " + err.Error()))
			return
		}
		hl.storage.AddGauge(mName, val)
	case "counter":
		val, err := utils.ParseCounter(mValue)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad counter value: " + err.Error()))
			return
		}
		hl.storage.AddCounter(mName, val)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unknown metric type: " + mType))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (hl *MMHandler) GetCurrentMetric(w http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, "type")
	mName := chi.URLParam(r, "name")
	switch mType {
	case "gauge":
		val, err := hl.storage.GetCurrentGauge(mName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Unknown metric: " + mName))
			return
		}
		w.Write([]byte(fmt.Sprintf("%v", val)))
	case "counter":
		val, err := hl.storage.GetCounter(mName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Unknown metric: " + mName))
			return
		}
		w.Write([]byte(fmt.Sprintf("%v", val)))
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unknown metric type: " + mType))
		return
	}
}

func (hl *MMHandler) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	templFile, err := os.ReadFile(hl.fsPath + "\\index.html")
	if err != nil {
		logger.Log(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	templ, err := template.New("allMetrics").Parse(string(templFile))
	if err != nil {
		logger.Log(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	gauges, err := hl.storage.GetCurrentGauges()
	if err != nil {
		logger.Log(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	counters, err := hl.storage.GetCounters()
	if err != nil {
		logger.Log(err.Error())
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
		logger.Log(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}
