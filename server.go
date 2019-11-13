package authplay

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/shipit/authplay/routes"
)

const (
	allLocal    = "0.0.0.0"
	defaultPort = 8888
)

// ServerOptions used to pass params to the Server
type ServerOptions struct {
	Port int
}

// Server is root daemon process
type Server interface {
	Start()
}

// NewDefaultServer creates default Server
func NewDefaultServer() (Server, error) {
	return NewServer(&ServerOptions{Port: defaultPort})
}

// NewServer creates Server
func NewServer(opts *ServerOptions) (Server, error) {
	addr := fmt.Sprintf("%s:%d", allLocal, opts.Port)

	router := mux.NewRouter()
	logger := handlers.LoggingHandler(os.Stdout, router)
	http := &http.Server{Handler: logger, Addr: addr}

	return &server{opts: opts, router: router, http: http}, nil
}

type server struct {
	opts   *ServerOptions
	http   *http.Server
	router *mux.Router
}

func (s *server) Start() {
	go func() {
		s.setupRoutes()
		if err := s.http.ListenAndServe(); err != nil {
			log.Panicf("%#v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.http.Shutdown(ctx)
}

func (s *server) setupRoutes() {
	s.router.HandleFunc("/", routes.Index)
}
