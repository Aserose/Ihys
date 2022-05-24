package server

import (
	"IhysBestowal/internal/config"
	"context"
	"net/http"
	"time"
)

type server struct {
	server *http.Server
}

func NewServer(cfg config.Server, handler http.Handler) *server {
	return &server{
		server: &http.Server{
			Addr:         ":" + cfg.Port,
			Handler:      handler,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

func (s server) Run() error { return s.server.ListenAndServe() }

func (s server) Shutdown(ctx context.Context) error { return s.server.Shutdown(ctx) }
