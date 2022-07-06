package app

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/repository"
	"IhysBestowal/internal/server"
	"IhysBestowal/internal/server/handler"
	"IhysBestowal/internal/service"
	"IhysBestowal/pkg/customLogger"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	log := customLogger.NewLogger()

	cfg := config.NewCfg(log)

	repo := repository.NewRepository(log, cfg.Repository)

	services := service.NewService(log, cfg.Service, repo)

	handlers := handler.NewHandler(log, cfg.Handler, services)

	srv := server.NewServer(cfg.Server, handlers.SetupRoutes())

	exit := newExit()

	go func() {
		if err := srv.Run(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(log.CallInfoStr(), err.Error())
			}
		}
	}()

	<-exit

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		services.Close()
		repo.Close()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(log.CallInfoStr(), err.Error())
	}

}

func newExit() chan os.Signal {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	return exit
}
