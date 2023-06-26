package server

import (
	"net/http"
	"os"
	"strings"

	"github.com/90amper/metmon/internal/config"
	"github.com/90amper/metmon/internal/logger"
	"github.com/90amper/metmon/internal/server/handlers"
	"github.com/90amper/metmon/internal/storage"
	"github.com/go-chi/chi/v5"

	mdw "github.com/90amper/metmon/internal/middleware"
	// "go.uber.org/zap"
)

type Server struct {
	Storage storage.Storager
	Router  *chi.Mux
	Handler *handlers.MMHandler
	FsPath  string
}

func (s *Server) NewRouter() *chi.Mux {
	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	r.Use(mdw.Logger)
	FileServer(r, "/html", http.Dir(s.FsPath))
	r.Get("/", s.Handler.GetAllMetrics)
	r.Route("/value", func(r chi.Router) {
		r.Post("/", s.Handler.GetCurrentJsonMetric)
		r.Route("/{type}", func(r chi.Router) {
			r.Get("/{name}", s.Handler.GetCurrentMetric)
		})
	})
	r.Route("/update", func(r chi.Router) {
		r.Post("/", s.Handler.ReceiveJsonMetrics)
		r.Route("/{type}", func(r chi.Router) {
			r.Route("/{name}", func(r chi.Router) {
				r.Post("/{value}", s.Handler.ReceiveMetrics)
			})
		})
	})
	// r.Post("/value", s.Handler.GetCurrentMetric)
	// r.Post("/update", s.Handler.ReceiveMetrics)
	return r
}

func NewServer() (srv *Server, err error) {
	srv = &Server{}
	srv.Storage = storage.NewStorage()
	wdPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	srv.FsPath = wdPath + "\\..\\..\\internal\\server\\html"
	srv.Handler, err = handlers.NewHandler(srv.Storage, srv.FsPath)
	srv.Router = srv.NewRouter()
	if err != nil {
		return nil, err
	}
	return srv, nil
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
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

func Run() (err error) {
	srv, err := NewServer()
	if err != nil {
		return err
	}
	logger.Log("Starting server at " + config.Config.ServerURL)
	err = http.ListenAndServe(config.Config.ServerURL, srv.Router)
	if err != nil {
		return (err)
	}
	return nil
}
