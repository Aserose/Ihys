package tg

import (
	"IhysBestowal/pkg/customLogger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goccy/go-json"
	"io"
)

type decoder struct {
	log customLogger.Logger
}

func newDecoder(log customLogger.Logger) decoder {
	return decoder{
		log: log,
	}
}

func (d decoder) parse(reqBody io.ReadCloser) tgbotapi.Update {
	incoming := tgbotapi.Update{}

	body, err := io.ReadAll(reqBody)
	if err != nil {
		d.log.Error(d.log.CallInfoStr(), err.Error())
	}

	err = json.Unmarshal(body, &incoming)
	if err != nil {
		d.log.Error(d.log.CallInfoStr(), err.Error())
	}

	return incoming
}
