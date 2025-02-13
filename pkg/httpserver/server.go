package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cutlery47/posts/config"
)

type Server struct {
	server          *http.Server
	shutDownTimeout time.Duration
}

func New(conf config.HTTPServer, handler http.Handler) *Server {
	std := &http.Server{
		Handler:      handler,
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
		Addr:         fmt.Sprintf("%v:%v", conf.BindAddress, conf.BindPort),
	}

	s := &Server{
		server:          std,
		shutDownTimeout: conf.ShutdownTimeout,
	}

	return s
}

func (s *Server) Run(errChan <-chan error) error {
	log.Println("[HTTPSERVER] listening on:", s.server.Addr)

	go func() {
		err := s.server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Println("[HTTPSERVER] http server error:", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		return s.shutdown(err)
	case <-sigChan:
		return s.shutdown(nil)
	}

}

func (s *Server) shutdown(e error) error {
	log.Println("[SHUTDOWN] http server shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), s.shutDownTimeout)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	return e
}
