package tg

import (
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/service"
	"IhysBestowal/pkg/customLogger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"
)

const (
	cmdMainMenu = `mainMenu`
	cmdAuthVk   = `authVk`
	cmdStart    = `start`
	cmdSearch   = `search`

	dlt = `delete`

	initMsg = ` ` // TODO
)

type Handler struct {
	service service.Service
	exe     dto.ExecCmd
	log     customLogger.Logger
	picker
	decoder
}

func New(log customLogger.Logger, service service.Service) Handler {
	return Handler{
		service: service,
		exe:     newExe(service),
		log:     log,
		decoder: newDecoder(log),
	}
}

func (h Handler) Webhook(w http.ResponseWriter, r *http.Request) {
	incoming := h.read(r.Body)

	switch incoming.Message {
	case nil:
		go h.exec(incoming)

	default:
		if incoming.Message.Command() != cmdStart {
			h.deleteMsg(incoming)
		}
		if incoming.Message.IsCommand() {
			switch incoming.Message.Command() {

			case cmdMainMenu:
				h.main(incoming)

			case cmdAuthVk:
				h.authVk(incoming)

			case cmdStart:
				// TODO

			case cmdSearch:
				h.search(incoming)

			}
		}
	}
}

func (h Handler) exec(incoming tgbotapi.Update) {
	callbackData, chatId, msgId := h.exeArgs(incoming)
	if h.exe[callbackData] == nil {

		//TODO

	} else {
		h.exe[callbackData](dto.Response{
			TGUser: dto.TGUser{
				ChatId: chatId,
			},
			MsgId:        msgId,
			MsgText:      incoming.CallbackQuery.Message.Text,
			CallbackData: callbackData,
			ExecCmd:      h.exe,
		})
	}
}

func (h Handler) init(incoming tgbotapi.Update) {
	_, chatId := h.userAndChatIDs(incoming)

	h.service.TG.Send(tgbotapi.NewMessage(chatId, initMsg))
}

func (h Handler) search(incoming tgbotapi.Update) {
	userId, chatId, query := h.cmdArgs(incoming)

	if len(query) != 0 {
		h.find(userId, chatId, query)
	} else {
		h.random(incoming)
	}
}

func (h Handler) random(incoming tgbotapi.Update) {
	userId, chatId := h.userAndChatIDs(incoming)

	h.service.Menu.Random(dto.Response{
		TGUser: dto.TGUser{
			UserId: userId,
			ChatId: chatId,
		},
		MsgId:   0,
		ExecCmd: h.exe,
	})
}

func (h Handler) authVk(incoming tgbotapi.Update) {
	userId, chatId, key := h.cmdArgs(incoming)

	// TODO
	//if !h.service.Auth.IsValidToken(key) {
	//	h.service.TG.Send(tgbotapi.NewMessage(chatId, "invalid"))
	//	return
	//}

	h.service.Auth.Vk().Create(dto.TGUser{
		UserId: userId,
		ChatId: chatId,
	}, key)

	h.main(incoming)
}

func (h Handler) main(incoming tgbotapi.Update) {
	userId, chatId := h.userAndChatIDs(incoming)

	h.service.Menu.Main(dto.Response{
		TGUser: dto.TGUser{
			UserId: userId,
			ChatId: chatId,
		},
		MsgId:   0,
		ExecCmd: h.exe,
	})
}

func (h Handler) find(userId, chatId int64, query string) {
	h.service.Menu.Find(
		dto.Response{
			TGUser: dto.TGUser{
				UserId: userId,
				ChatId: chatId,
			},
			MsgId:   0,
			ExecCmd: h.exe,
		}, query)
}

func (h Handler) deleteMsg(incoming tgbotapi.Update) {
	chatId, msgId := h.userAndMsgIDs(incoming)
	h.exe[dlt](dto.Response{
		TGUser: dto.TGUser{
			ChatId: chatId,
		},
		MsgId: msgId,
	})
}
