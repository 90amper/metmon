package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/90amper/metmon/internal/models"
)

func (hl *MMHandler) ReceiveJSONMetrics(w http.ResponseWriter, r *http.Request) {
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
		if mValue == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		hl.storage.AddGauge(mName, models.Gauge(*mValue))
		curVal, err := hl.storage.GetCurrentGauge(mName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		val := float64(curVal)
		resp.Value = &val
	case "counter":
		if mDelta == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		hl.storage.AddCounter(mName, models.Counter(*mDelta))
		curVal, err := hl.storage.GetCounter(mName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		val := int64(curVal)
		resp.Delta = &val
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unknown metric type: " + mType))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (hl *MMHandler) GetCurrentJSONMetric(w http.ResponseWriter, r *http.Request) {
	var req models.Metric
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mType := req.MType
	mName := req.ID
	resp := models.Metric{
		ID:    mName,
		MType: mType,
	}
	switch mType {
	case "gauge":
		curVal, err := hl.storage.GetCurrentGauge(mName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		val := float64(curVal)
		resp.Value = &val
	case "counter":
		curVal, err := hl.storage.GetCounter(mName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		val := int64(curVal)
		resp.Delta = &val
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unknown metric type: " + mType))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
