package tg

import (
	"IhysBestowal/internal/config"
	tgApi "IhysBestowal/internal/service/webapi/tg/api"
	"IhysBestowal/internal/service/webapi/tg/menu"
	"IhysBestowal/pkg/customLogger"
)

type Telegram struct {
	tgApi.Api
	menu.Builder
}

func New(log customLogger.Logger, cfg config.Service) Telegram {
	api := tgApi.New(log, cfg.Telegram)

	return Telegram{
		Api:     api,
		Builder: menu.New(api),
	}
}

func (t Telegram) Menu() menu.Builder {
	return t.Builder
}
