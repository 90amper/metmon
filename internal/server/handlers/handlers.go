package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/90amper/metmon/internal/models"
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
	// w.Write([]byte("All metrics"))
	// templFile, err := os.ReadFile("../internal/server/static/temp/index1.html")
	// err := mime.AddExtensionType(".css", "text/css")
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }
	templPath, err := os.Getwd()
	templPath += "\\..\\..\\internal\\server\\html\\index.html"
	fmt.Println(templPath)
	templFile, err := os.ReadFile(templPath)
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

	// test := make(map[string]string)
	// test["one"] = "ONE"
	// test["two"] = "TWO"
	data := struct {
		Gauges   models.GaugeList
		Counters models.CounterStore
	}{
		Gauges:   gauges,
		Counters: counters,
	}

	err = templ.Execute(w, data)

}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
