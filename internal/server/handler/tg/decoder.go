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

func (d decoder) read(reqBody io.ReadCloser) tgbotapi.Update {
	incoming := tgbotapi.Update{}
	d.parseReqBody(reqBody, &incoming)
	return incoming
}

func (d decoder) parseReqBody(reqBody io.ReadCloser, v interface{}) {
	d.unmarshal(d.readBody(reqBody), v)
}

func (d decoder) readBody(reqBody io.ReadCloser) []byte {
	body, err := io.ReadAll(reqBody)
	if err != nil {
		d.log.Error(d.log.CallInfoStr(), err.Error())
	}

	return body
}

func (d decoder) unmarshal(data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	if err != nil {
		d.log.Error(d.log.CallInfoStr(), err.Error())
	}
}
