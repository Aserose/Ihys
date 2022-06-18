package handler

import (
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/service"
	"IhysBestowal/pkg/customLogger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goccy/go-json"
	"io"
	"net/http"
	"strings"
)

type tgHandler struct {
	service service.Service
	exe     dto.ExecCmd
	p       picker
	log     customLogger.Logger
}

func newTGHandler(log customLogger.Logger, service service.Service) tgHandler {
	return tgHandler{
		service: service,
		exe: map[string]dto.OnTappedFunc{
			"delete": func(p dto.Response) {
				service.SendMsg(tgbotapi.NewDeleteMessage(p.ChatId, p.MsgId))
			},
		},
		log: log,
	}
}

func (h tgHandler) mainWebhook(w http.ResponseWriter, r *http.Request) {
	update := h.getUpdate(h.log, r.Body)

	switch update.Message {
	case nil:
		h.execute(update)

	default:
		if update.Message.Command() != "start" {
			h.deleteMsg(update)
		}
		if update.Message.IsCommand() {
			switch update.Message.Command() {

			case "mainMenu":
				h.openMainMenu(update)

			case "authVk":
				h.authVk(update)

			case "start":
				// TODO

			case "search":
				h.search(update)

			}
		}
	}
}

func (h tgHandler) execute(incoming tgbotapi.Update) {
	callbackData, chatId, msgId := h.p.getExeArgs(incoming)
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

func (h tgHandler) search(incoming tgbotapi.Update) {
	userId, chatId, query := h.p.getCmdArgs(incoming)

	h.openSearchMenu(userId, chatId, query)
}

func (h tgHandler) authVk(incoming tgbotapi.Update) {
	userId, chatId, key := h.p.getCmdArgs(incoming)

	if !h.service.AuthService.IsValidToken(key) {
		h.service.TelegramService.SendMsg(tgbotapi.NewMessage(chatId, "invalid token"))
		return
	}

	h.service.AuthService.Vk().PutKey(dto.TGUser{
		UserId: userId,
		ChatId: chatId,
	}, key)

	h.openMainMenu(incoming)
}

func (h tgHandler) openMainMenu(incoming tgbotapi.Update) {
	userId, chatId := h.p.getUserAndChatIDs(incoming)

	h.service.TGMenu.Main(dto.Response{
		TGUser: dto.TGUser{
			UserId: userId,
			ChatId: chatId,
		},
		MsgId:   0,
		ExecCmd: h.exe,
	})
}

func (h tgHandler) openSearchMenu(userId, chatId int64, query string) {
	h.service.TGMenu.Search(
		dto.Response{
			TGUser: dto.TGUser{
				UserId: userId,
				ChatId: chatId,
			},
			MsgId:   0,
			ExecCmd: h.exe,
		}, query)
}

func (h tgHandler) deleteMsg(incoming tgbotapi.Update) {
	chatId, msgId := h.p.getUserAndMsgIDs(incoming)
	h.exe["delete"](dto.Response{
		TGUser: dto.TGUser{
			ChatId: chatId,
		},
		MsgId: msgId,
	})
}

func (h tgHandler) getUpdate(log customLogger.Logger, reqBody io.ReadCloser) tgbotapi.Update {
	incoming := tgbotapi.Update{}
	h.parseRequestBody(log, reqBody, &incoming)
	return incoming
}

func (h tgHandler) parseRequestBody(log customLogger.Logger, reqBody io.ReadCloser, v interface{}) {
	h.unmarshal(log, h.readBody(log, reqBody), v)
}

func (h tgHandler) readBody(log customLogger.Logger, reqBody io.ReadCloser) []byte {
	body, err := io.ReadAll(reqBody)
	if err != nil {
		log.Error(log.CallInfoStr(), err.Error())
	}

	return body
}

func (h tgHandler) unmarshal(log customLogger.Logger, data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	if err != nil {
		log.Error(log.CallInfoStr(), err.Error())
	}
}

// picker is responsible for retrieving certain parameters from the update response
type picker struct{}

func (p picker) getExeArgs(incoming tgbotapi.Update) (callbackData string, chatId int64, msgId int) {
	callbackData = incoming.CallbackQuery.Data
	chatId, msgId = p.getCallbackChatAndMsgIDs(incoming)

	return
}

func (p picker) getCmdArgs(incoming tgbotapi.Update) (userId, chatId int64, query string) {
	userId, chatId = p.getUserAndChatIDs(incoming)
	query = p.getCmdContent(incoming.Message.Text, incoming.Message.Command())

	return
}

func (p picker) getCallbackChatAndMsgIDs(incoming tgbotapi.Update) (chatId int64, msgId int) {
	chatId = incoming.CallbackQuery.Message.Chat.ID
	msgId = incoming.CallbackQuery.Message.MessageID

	return
}

func (p picker) getUserAndChatIDs(incoming tgbotapi.Update) (userId, chatId int64) {
	userId = incoming.SentFrom().ID
	chatId = incoming.Message.Chat.ID

	return
}

func (p picker) getCmdContent(rawMsgText, nameCmd string) string {
	result := strings.Split(rawMsgText, nameCmd+" ")
	if len(result) <= 1 {
		return " "
	}
	return strings.Split(rawMsgText, nameCmd+" ")[1]
}

func (p picker) getUserAndMsgIDs(incoming tgbotapi.Update) (userId int64, msgId int) {
	userId = incoming.SentFrom().ID
	msgId = incoming.Message.MessageID

	return
}
