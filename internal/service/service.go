package service

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/repository"
	"IhysBestowal/internal/service/auth"
	"IhysBestowal/internal/service/menu"
	"IhysBestowal/internal/service/webapi"
	"IhysBestowal/pkg/customLogger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Auth interface {
	Vk() auth.Key
}

type TG interface {
	Send(c tgbotapi.Chattable) tgbotapi.Message
}

type Menu interface {
	Main(p dto.Response)
	Find(p dto.Response, query string)
	Random(p dto.Response)
	Init(p dto.Response)
}

type Service struct {
	Auth
	TG
	Menu
	exit []func()
}

func New(log customLogger.Logger, cfg config.Service, repo repository.Repository) Service {
	a := auth.New(log, cfg.Auth, repo)
	wa := webapi.New(log, cfg, repo, a)

	return Service{
		Auth: a,
		TG:   wa.TG,
		Menu: menu.New(wa, repo, cfg.Keypads),
		exit: []func(){wa.Close},
	}
}

func (s Service) Close() {
	for _, e := range s.exit {
		e()
	}
}
