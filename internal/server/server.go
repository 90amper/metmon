package server

import (
	"net/http"
	"time"

	"github.com/90amper/metmon/internal/config"
	"github.com/90amper/metmon/internal/logger"
	"github.com/90amper/metmon/internal/server/handlers"
	"github.com/90amper/metmon/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	timeout time.Duration
	deep    int64
	Storage storage.Storager
	Router  *chi.Mux
	Wrapper *handlers.Wrapper
}

func (s *Server) NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", s.Wrapper.GetAllMetrics)
	r.Route("/value", func(r chi.Router) {
		r.Route("/{type}", func(r chi.Router) {
			r.Get("/{name}", s.Wrapper.GetCurrentMetric)
		})
	})
	r.Route("/update", func(r chi.Router) {
		r.Route("/{type}", func(r chi.Router) {
			r.Route("/{name}", func(r chi.Router) {
				r.Post("/{value}", s.Wrapper.ReceiveMetrics)
			})
		})
	})
	// r.Get("/value/{type}/{name}", s.Wrapper.GetCurrentMetric)
	// r.Post("/update/{type}/{name}/{value}", s.Wrapper.ReceiveMetrics)
	return r
}

func NewServer() (srv *Server, err error) {
	srv = &Server{}
	srv.Storage = storage.NewStorage()
	srv.Wrapper, err = handlers.NewWrapper(srv.Storage)
	srv.Router = srv.NewRouter()

	if err != nil {
		return nil, err
	}
	return srv, nil
}

func Run() (err error) {
	srv, err := NewServer()
	if err != nil {
		return err
	}
	logger.Log("Starting server at " + config.CmdFlags.ServerUrl)
	err = http.ListenAndServe(config.CmdFlags.ServerUrl, srv.Router)
	if err != nil {
		return (err)
	}
	return nil
}

// func collectorHandler(w http.ResponseWriter, r *http.Request) {
// 	// if r.Method == "POST" && r.Header.Get("Content-Type") == "text/plain" {
// 	fmt.Println(r.Method, r.URL.Path)
// 	if r.Method == "POST" {
// 		// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
// 		path := strings.Split(r.URL.Path, "/")

// 		if len(path) < 5 {
// 			w.WriteHeader(http.StatusNotFound)
// 			return
// 		}

// 		if path[2] == "counter" {
// 			vali, err := strconv.ParseInt(path[4], 10, 64)
// 			if err != nil {
// 				w.WriteHeader(http.StatusBadRequest)
// 				return
// 			}
// 			fmt.Printf("%+v\n", vali)
// 			// MemStor.counter[path[3]] += vali
// 			// w.Write([]byte(fmt.Sprintf("%+v\r\n", MemStor)))
// 			return
// 		} else if path[2] == "gauge" {
// 			valf, err := strconv.ParseFloat(path[4], 64)
// 			if err != nil {
// 				w.WriteHeader(http.StatusBadRequest)
// 				return
// 			}
// 			// MemStor.gauge[path[3]] = valf
// 			// w.Write([]byte(fmt.Sprintf("%+v\r\n", MemStor)))
// 			fmt.Printf("%+v\n", valf)
// 			return
// 		} else {
// 			// println(path)
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
// 		// fmt.Printf("%+v\r\n", MemStor)

// 		// test := strings.Join(path, " | ")

// 		// w.WriteHeader(http.StatusInternalServerError)
// 		// w.Write([]byte("test: " + test))
// 		// return
// 	} else {
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 		return
// 	}
// }
