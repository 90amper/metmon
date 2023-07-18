package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/90amper/metmon/internal/logger"
	"github.com/90amper/metmon/internal/models"
)

func (hl *MMHandler) ReceiveJSONMetric(w http.ResponseWriter, r *http.Request) {
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
		err := hl.storage.AddGauge(mName, models.Gauge(*mValue))
		if err != nil {
			logger.Error(fmt.Errorf("add gauge failed: %w", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		curVal, err := hl.storage.GetCurrentGauge(mName)
		if err != nil {

			logger.Error(fmt.Errorf("read gauge failed: %w", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		val := float64(curVal)
		logger.Trace("Current %s: %v", mName, val)
		resp.Value = &val
	case "counter":
		if mDelta == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err := hl.storage.AddCounter(mName, models.Counter(*mDelta))
		if err != nil {
			logger.Error(fmt.Errorf("add counter failed: %w", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		curVal, err := hl.storage.GetCounter(mName)
		if err != nil {
			// logger.Log(err.Error())
			logger.Error(fmt.Errorf("read counter failed: %w", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		val := int64(curVal)
		logger.Debug("Current %s: %v", mName, val)
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

func (hl *MMHandler) ReceiveJSONMetrics(w http.ResponseWriter, r *http.Request) {
	var errs []error
	// var resp []models.Metric
	var reqArr []models.Metric
	err := json.NewDecoder(r.Body).Decode(&reqArr)
	if err != nil {
		logger.Error(fmt.Errorf("can't decode metric: %w", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		// mock, _ := json.Marshal([]string{})
		// http.Error(w, "[]", http.StatusBadRequest)
		// w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(http.StatusBadRequest)
		// json.NewEncoder(w).Encode([]string{})
		return
	}

	for _, req := range reqArr {

		mType := req.MType
		mName := req.ID
		mValue := req.Value
		mDelta := req.Delta

		if mName == "" {
			// w.WriteHeader(http.StatusNotFound)
			// w.Write([]byte("Metric name not specified!"))
			// return
			errs = append(errs, fmt.Errorf("metric name not specified"))
			continue
		}
		switch mType {
		case "gauge":
			if mValue == nil {
				// w.WriteHeader(http.StatusBadRequest)
				// return
				errs = append(errs, fmt.Errorf("empty gauge value"))
				continue
			}
			err := hl.storage.AddGauge(mName, models.Gauge(*mValue))
			if err != nil {
				logger.Error(fmt.Errorf("add gauge failed: %w", err))
				errs = append(errs, err)
				continue
				// w.WriteHeader(http.StatusInternalServerError)
				// return
			}

			// curVal, err := hl.storage.GetCurrentGauge(mName)
			// if err != nil {

			// 	logger.Error(fmt.Errorf("read gauge failed: %w", err))
			// 	w.WriteHeader(http.StatusInternalServerError)
			// 	return
			// }
			// val := float64(curVal)
			// logger.Trace("Current %s: %v", mName, val)
			// // resp.Value = &val
			// resp = append(resp, models.Metric{
			// 	ID:    mName,
			// 	MType: mType,
			// 	Value: &val,
			// })
		case "counter":
			if mDelta == nil {
				// w.WriteHeader(http.StatusBadRequest)
				// return
				errs = append(errs, fmt.Errorf("empty counter value"))
				continue
			}
			err := hl.storage.AddCounter(mName, models.Counter(*mDelta))
			if err != nil {
				logger.Error(fmt.Errorf("add counter failed: %w", err))
				errs = append(errs, err)
				continue
				// w.WriteHeader(http.StatusInternalServerError)
				// return
			}
			// curVal, err := hl.storage.GetCounter(mName)
			// if err != nil {
			// 	// logger.Log(err.Error())
			// 	logger.Error(fmt.Errorf("read counter failed: %w", err))
			// 	w.WriteHeader(http.StatusInternalServerError)
			// 	return
			// }
			// val := int64(curVal)
			// logger.Trace("Current %s: %v", mName, val)
			// // resp.Delta = &val
			// resp = append(resp, models.Metric{
			// 	ID:    mName,
			// 	MType: mType,
			// 	Delta: &val,
			// })
		default:
			// w.WriteHeader(http.StatusBadRequest)
			// w.Write([]byte("Unknown metric type: " + mType))
			errs = append(errs, fmt.Errorf("unknown metric type: %s", mType))
			continue
		}
	}
	// w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprint(errors.Join(errs...))))

	// json.NewEncoder(w).Encode(resp)

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
