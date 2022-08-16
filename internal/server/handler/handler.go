package handler

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/server/handler/tg"
	"IhysBestowal/internal/service"
	"IhysBestowal/pkg/customLogger"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	api config.Api
	tg  tg.Handler
}

func New(log customLogger.Logger, cfg config.Handler, service service.Service) Handler {
	return Handler{
		api: cfg.Api,
		tg:  tg.New(log, service),
	}
}

func (h Handler) SetupRoutes() *chi.Mux {
	router := chi.NewRouter()

	telegram := router.Group(func(r chi.Router) {})
	{
		telegram.Post(h.api.Telegram, h.tg.Webhook)
	}

	return router
}
