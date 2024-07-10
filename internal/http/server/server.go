package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Lanworm/OTUS_GO/final_project/internal/config"
	"github.com/Lanworm/OTUS_GO/final_project/internal/logger"
)

type Server struct {
	logger *logger.Logger
	conf   config.ServerHTTPConf
	srv    *http.Server
	mux    *http.ServeMux
	done   bool
}

func NewHTTPServer(
	logger *logger.Logger,
	conf config.ServerHTTPConf,
) *Server {
	if conf.Protocol == "" {
		conf.Protocol = "tcp4"
	}

	return &Server{
		logger: logger,
		conf:   conf,
		mux:    http.NewServeMux(),
	}
}

func (s *Server) Start(_ context.Context) error {
	if s.srv != nil {
		return errors.New("http server already started")
	}

	lw := NewLogMiddleware(s.logger)
	lc := NewRecoveryMiddleware(s.logger)

	s.srv = &http.Server{
		Addr:              s.conf.GetFullAddress(),
		Handler:           lc.Wrap(lw.Wrap(s.mux)),
		TLSConfig:         nil,
		ReadTimeout:       s.conf.Timeout,
		ReadHeaderTimeout: s.conf.Timeout,
		WriteTimeout:      s.conf.Timeout,
		IdleTimeout:       s.conf.Timeout,
		MaxHeaderBytes:    1 << 10,
	}

	err := s.srv.ListenAndServe()
	if err != nil && !(errors.Is(err, http.ErrServerClosed) && s.done) {
		return fmt.Errorf("http listen and serve at {%s}: %w", s.conf.GetFullAddress(), err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.done = true

	return s.srv.Shutdown(ctx)
}

func (s *Server) AddRoute(route string, handlerFunc http.HandlerFunc) {
	s.mux.HandleFunc(route, handlerFunc)
}
