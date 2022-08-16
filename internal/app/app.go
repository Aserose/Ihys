package app

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/repository"
	"IhysBestowal/internal/server"
	"IhysBestowal/internal/server/handler"
	"IhysBestowal/internal/service"
	"IhysBestowal/pkg/customLogger"
	"context"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() { rand.Seed(time.Now().UnixNano()) }

func Run() {
	log := customLogger.NewLogger()

	cfg := config.New(log)

	repo := repository.New(log, cfg.Repository)

	services := service.New(log, cfg.Service, repo)

	handlers := handler.New(log, cfg.Handler, services)

	srv := server.New(cfg.Server, handlers.SetupRoutes())

	go func() {
		if err := srv.Run(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(log.CallInfoStr(), err.Error())
			}
		}
	}()

	<-exit()

	services.Close()
	repo.Close()
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal(log.CallInfoStr(), err.Error())
	}

}

func exit() chan os.Signal {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	return exit
}
