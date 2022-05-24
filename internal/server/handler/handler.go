package handler

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/service"
	"IhysBestowal/pkg/customLogger"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service service.Service
	api     config.Api
	tg      tgHandler
	log     customLogger.Logger
}

func NewHandler(log customLogger.Logger, cfg config.Handler, service service.Service) Handler {
	return Handler{
		service: service,
		api:     cfg.Api,
		tg:      newTGHandler(log, service),
		log:     log,
	}
}

func (h Handler) SetupRoutes() *chi.Mux {
	router := chi.NewRouter()

	telegram := router.Group(func(r chi.Router) {})
	{
		telegram.Post(h.api.Telegram, h.tg.mainWebhook)
	}

	return router
}
