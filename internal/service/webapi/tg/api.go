package tg

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/pkg/customLogger"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Api interface {
	Send(chattable tgbotapi.Chattable) tgbotapi.Message
}

type api struct {
	*tgbotapi.BotAPI
	log customLogger.Logger
}

func newApi(log customLogger.Logger, cfg config.Telegram) Api {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Fatal(log.CallInfoStr(), err.Error())
	}

	wh, err := tgbotapi.NewWebhook(cfg.WebhookLink)
	if err != nil {
		log.Print(err.Error())
	}

	_, err = bot.Request(wh)
	if err != nil {
		log.Fatal(log.CallInfoStr(), err.Error())
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(log.CallInfoStr(), err.Error())
	}

	if info.LastErrorDate != 0 {
		log.Error(log.CallInfoStr(), fmt.Sprintf("Telegram callback failed: %s", info.LastErrorMessage))
	}

	bot.Debug = true

	return api{
		BotAPI: bot,
		log:    log,
	}
}

func (t api) Send(chattable tgbotapi.Chattable) tgbotapi.Message {
	a, err := t.BotAPI.Send(chattable)
	if err != nil {
		t.log.Error(t.log.CallInfoStr(), err.Error())
	}
	return a
}
