package tg

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/service/webapi/tg/api"
	"IhysBestowal/internal/service/webapi/tg/menu"
	"IhysBestowal/pkg/customLogger"
)

type Telegram struct {
	api.Api
	menu.Builder
}

func New(log customLogger.Logger, cfg config.Service) Telegram {
	a := api.New(log, cfg.Telegram)

	return Telegram{
		Api:     a,
		Builder: menu.New(a),
	}
}

func (t Telegram) Menu() menu.Builder {
	return t.Builder
}
