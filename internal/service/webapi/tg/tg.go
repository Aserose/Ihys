package tg

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/pkg/customLogger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ITelegram interface {
	Send(chattable tgbotapi.Chattable) tgbotapi.Message
	NewMenuBuilder() TGMenu
}

type tg struct {
	Api
	TGMenu
}

func NewTg(log customLogger.Logger, cfg config.Service) ITelegram {
	api := newApi(log, cfg.Telegram)

	return tg{
		Api:    api,
		TGMenu: newMenuBuilder(api, cfg.Menu),
	}
}

func (t tg) NewMenuBuilder() TGMenu {
	return t.TGMenu
}
