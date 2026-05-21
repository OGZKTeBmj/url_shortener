package httpserver

import (
	"context"
	"errors"
	"net/http"
	"time"
)

type Config struct {
	ShutdownTimeout time.Duration `default:"10s" envconfig:"HTTP_SHUTDOWN_TIMEOUT"`
	Port            string        `default:":8080" envconfig:"HTTP_PORT"`
}

type Server struct {
	server *http.Server
	cfg    Config
}

func New(handler http.Handler, cfg Config) *Server {
	return &Server{
		server: &http.Server{
			Addr:    ":" + cfg.Port,
			Handler: handler,
		},
		cfg: cfg,
	}
}

func (s *Server) Run() error {
	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		s.cfg.ShutdownTimeout,
	)
	defer cancel()

	return s.server.Shutdown(ctx)
}
