package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func Run() {
	r := chi.NewRouter()
	r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("chi"))
	})

	err := http.ListenAndServe(`:8080`, r)

	if err != nil {
		panic(err)
	}
}

func collectorHandler(w http.ResponseWriter, r *http.Request) {
	// if r.Method == "POST" && r.Header.Get("Content-Type") == "text/plain" {
	fmt.Println(r.Method, r.URL.Path)
	if r.Method == "POST" {
		// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
		path := strings.Split(r.URL.Path, "/")

		if len(path) < 5 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if path[2] == "counter" {
			vali, err := strconv.ParseInt(path[4], 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			fmt.Printf("%+v\n", vali)
			// MemStor.counter[path[3]] += vali
			// w.Write([]byte(fmt.Sprintf("%+v\r\n", MemStor)))
			return
		} else if path[2] == "gauge" {
			valf, err := strconv.ParseFloat(path[4], 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			// MemStor.gauge[path[3]] = valf
			// w.Write([]byte(fmt.Sprintf("%+v\r\n", MemStor)))
			fmt.Printf("%+v\n", valf)
			return
		} else {
			// println(path)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// fmt.Printf("%+v\r\n", MemStor)

		// test := strings.Join(path, " | ")

		// w.WriteHeader(http.StatusInternalServerError)
		// w.Write([]byte("test: " + test))
		// return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
