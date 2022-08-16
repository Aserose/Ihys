package server

import (
	"IhysBestowal/internal/config"
	"context"
	"net/http"
)

type server struct {
	server *http.Server
}

func New(cfg config.Server, handler http.Handler) *server {
	return &server{
		server: &http.Server{
			Addr:    ":" + cfg.Port,
			Handler: handler,
		},
	}
}

func (s server) Run() error                         { return s.server.ListenAndServe() }
func (s server) Shutdown(ctx context.Context) error { return s.server.Shutdown(ctx) }
