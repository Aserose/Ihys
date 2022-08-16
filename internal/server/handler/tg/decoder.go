package tg

import (
	"IhysBestowal/pkg/customLogger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goccy/go-json"
	"io"
)

type decoder struct{}

func (d decoder) update(log customLogger.Logger, reqBody io.ReadCloser) tgbotapi.Update {
	incoming := tgbotapi.Update{}
	d.parseRequestBody(log, reqBody, &incoming)
	return incoming
}

func (d decoder) parseRequestBody(log customLogger.Logger, reqBody io.ReadCloser, v interface{}) {
	d.unmarshal(log, d.readBody(log, reqBody), v)
}

func (d decoder) readBody(log customLogger.Logger, reqBody io.ReadCloser) []byte {
	body, err := io.ReadAll(reqBody)
	if err != nil {
		log.Error(log.CallInfoStr(), err.Error())
	}

	return body
}

func (d decoder) unmarshal(log customLogger.Logger, data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	if err != nil {
		log.Error(log.CallInfoStr(), err.Error())
	}
}
