package app

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/repository"
	"IhysBestowal/internal/server"
	"IhysBestowal/internal/server/handler"
	"IhysBestowal/internal/service"
	"IhysBestowal/pkg/customLogger"
)

func Run() {
	log := customLogger.NewLogger()

	cfg := config.NewCfg(log)

	repo := repository.NewRepository(log, cfg.Repository)

	services := service.NewService(log, cfg.Service, repo)

	handlers := handler.NewHandler(log, cfg.Handler, services)

	err := server.NewServer(cfg.Server, handlers.SetupRoutes()).Run(); if err != nil {
		log.Fatal(log.CallInfoStr(), err.Error())
	}
}
