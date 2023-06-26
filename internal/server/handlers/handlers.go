package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/90amper/metmon/internal/models"
	"github.com/90amper/metmon/internal/storage"
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
	var req models.Metric
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mType := req.MType
	mName := req.ID
	mValue := req.Value
	mDelta := req.Delta

	if mName == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Metric name not specified!"))
		return
	}
	resp := models.Metric{
		ID:    mName,
		MType: mType,
	}
	switch mType {
	case "gauge":
		hl.storage.AddGauge(mName, models.Gauge(*mValue))
		curVal, err := hl.storage.GetCurrentGauge(mName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		val := float64(curVal)
		resp.Value = &val
	case "counter":
		hl.storage.AddCounter(mName, models.Counter(*mDelta))
		curVal, err := hl.storage.GetCounter(mName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		val := int64(curVal)
		resp.Delta = &val
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unknown metric type: " + mType))
	}
	json.NewEncoder(w).Encode(resp)
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
	}

}

func (hl *MMHandler) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	templFile, err := os.ReadFile(hl.fsPath + "\\index.html")
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

	gauges, err := hl.storage.GetCurrentGauges()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	counters, err := hl.storage.GetCounters()
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
