package http

import (
	"context"
	"github.com/dnozdrin/detask/internal/app/log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type server struct {
	http *http.Server
	log  log.Logger
}

const (
	writeTimeout = 15 * time.Second
	readTimeout = 15 * time.Second
)

// NewServer will create a new instance of the web server
func NewServer(handler http.Handler, log log.Logger) *server {
	return &server{
		log: log,
		http: &http.Server{
			Handler:      handler,
			WriteTimeout: writeTimeout,
			ReadTimeout:  readTimeout,
		},
	}
}

// Start will start the web server
func (s *server) Start(addr string) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	go func() {
		<-quit
		s.shutdown()
	}()

	s.log.Infof("starting http server on address %s", addr)
	s.http.Addr = addr

	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.log.Error("http: server: listen and server")
		s.log.Error(err)
	}
}

func (s *server) shutdown() {
	err := s.http.Shutdown(context.Background())
	if err != nil {
		s.log.Error("http: server: shutdown")
		s.log.Error(err)
	}
}
