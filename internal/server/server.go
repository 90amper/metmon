package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/90amper/metmon/internal/logger"
	"github.com/90amper/metmon/internal/server/config"
	"github.com/90amper/metmon/internal/server/handlers"
	"github.com/90amper/metmon/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	mdw "github.com/90amper/metmon/internal/middleware"
	// "go.uber.org/zap"
)

type Server struct {
	Storage storage.Storager
	Router  *chi.Mux
	Handler *handlers.MMHandler
	FsPath  string
	Ctx     context.Context
}

func (s *Server) NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		// AllowedOrigins: []string{"https://*", "http://*"},
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(mdw.GzipMiddleware)
	// r.Use(middleware.Logger)
	r.Use(mdw.Logger)
	FileServer(r, "/html", http.Dir(s.FsPath))
	r.Get("/", s.Handler.GetAllMetrics)
	r.Route("/value", func(r chi.Router) {
		r.Post("/", s.Handler.GetCurrentJSONMetric)
		r.Route("/{type}", func(r chi.Router) {
			r.Get("/{name}", s.Handler.GetCurrentMetric)
		})
	})
	r.Route("/update", func(r chi.Router) {
		r.Post("/", s.Handler.ReceiveJSONMetrics)
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
	srv.Ctx = context.Background()
	srv.Storage = storage.NewStorage(&config.Config)

	srv.FsPath = strings.Join([]string{config.Config.ProjPath, "internal", "server", "html", ""}, config.Config.PathSeparator)

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
	var signals = []os.Signal{
		// os.Interrupt,
		syscall.SIGINT,
		syscall.SIGQUIT,
		// syscall.SIGABRT,
		// syscall.SIGKILL,
		// syscall.SIGTERM,
		// syscall.SIGSTOP,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, signals...)

	// rootCtx := context.Background()
	// taskCtx, cancelFn := context.WithCancel(rootCtx)

	srv, err := NewServer()
	if err != nil {
		return err
	}
	if config.Config.Restore {
		err = srv.Storage.LoadFromFile()
		if err != nil {
			logger.Log(err.Error())
		}
	}

	fmt.Printf("%v Starting server at %v\n", time.Now().Format(time.RFC3339), config.Config.ServerURL)
	go func() {
		err = http.ListenAndServe(config.Config.ServerURL, srv.Router)
		if err != nil {
			logger.Log(err.Error())
			panic(err.Error())
		}
	}()

	if config.Config.FileStoragePath != "" {
		go srv.Storage.Dumper()
	}

	<-shutdown
	// cancelFn()
	fmt.Printf("%v Shutdown MetMon server ... ", time.Now().Format(time.RFC3339))
	srv.Storage.SaveToFile()
	fmt.Printf("done\n")
	return nil
}
