package handler

import (
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/service"
	"IhysBestowal/pkg/customLogger"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"net/http"
	"strings"
)

type tgHandler struct {
	service service.Service
	exe     map[string]func(chatId int64, msgId int)
	log     customLogger.Logger
}

func newTGHandler(log customLogger.Logger, service service.Service) tgHandler {
	exes := map[string]func(chatId int64, msgId int){
		"delete": func(chatId int64, msgId int) {
			service.SendMsg(tgbotapi.NewDeleteMessage(chatId, msgId))
		},
	}

	return tgHandler{
		service: service,
		exe:     exes,
		log:     log,
	}
}

func (h tgHandler) mainWebhook(w http.ResponseWriter, r *http.Request) {
	incoming := h.getUpdate(h.log, r.Body)

	switch incoming.Message {
	case nil:
		h.execute(h.getExeArgs(incoming))

	default:
		if incoming.Message.IsCommand() {
			switch incoming.Message.Command() {

			case "mainMenu":
				h.openMainMenu(h.getUserAndChatIDs(incoming))

			case "authVk":
				h.authVk(h.getCmdArgs(incoming))

			case "start":
				// TODO

			case "search":
				h.search(h.getCmdArgs(incoming))

			}
		}
	}
}

func (h tgHandler) execute(callbackData string, chatId int64, msgId int) {
	if h.exe[callbackData] == nil {

		//TODO

	} else {
		h.exe[callbackData](chatId, msgId)
	}
}

func (h tgHandler) search(userId, chatId int64, msgId int, query string) {
	h.openSearchMenu(userId, chatId, query)
	h.deleteMsg(chatId, msgId)
}

func (h tgHandler) authVk(userId, chatId int64, msgId int, key string) {
	if !h.service.AuthService.IsValidToken(key) {
		h.service.TelegramService.SendMsg(tgbotapi.NewMessage(chatId, "invalid token"))
		return
	}
	h.putKey(userId, chatId, key)
	h.deleteMsg(userId, msgId)
	h.openMainMenu(userId, chatId)
}

func (h tgHandler) openMainMenu(userId, chatId int64) {
	h.service.TGMenu.Main(dto.Executor{
		TGUser: dto.TGUser{
			UserId: userId,
			ChatID: chatId,
		},
		ExecCmd: h.exe,
	}, 0, nil)
}

func (h tgHandler) openSearchMenu(userId, chatId int64, query string) {
	h.service.TGMenu.Search(dto.Executor{
		TGUser: dto.TGUser{
			UserId: userId,
			ChatID: chatId,
		},
		ExecCmd: h.exe,
	}, 0, query)
}

func (h tgHandler) getCmdContent(rawMsgText, nameCmd string) string {
	result := strings.Split(rawMsgText, nameCmd+" ")
	if len(result) <= 1 {
		return " "
	}
	return strings.Split(rawMsgText, nameCmd+" ")[1]
}

func (h tgHandler) deleteMsg(chatId int64, msgId int) {
	h.exe["delete"](chatId, msgId)
}

func (h tgHandler) putKey(userId, chatId int64, key string) {
	h.service.AuthService.Vk().PutKey(dto.TGUser{
		UserId: userId,
		ChatID: chatId,
	}, key)
}

func (h tgHandler) getExeArgs(incoming tgbotapi.Update) (callbackData string, chatId int64, msgId int) {
	callbackData = incoming.CallbackQuery.Data
	chatId, msgId = h.getCallbackChatAndMsgIDs(incoming)

	return
}

func (h tgHandler) getCmdArgs(incoming tgbotapi.Update) (userId, chatId int64, msgId int, query string) {
	userId, chatId = h.getUserAndChatIDs(incoming)
	msgId = incoming.Message.MessageID
	query = h.getCmdContent(incoming.Message.Text, incoming.Message.Command())

	return
}

func (h tgHandler) getCallbackChatAndMsgIDs(incoming tgbotapi.Update) (chatId int64, msgId int) {
	chatId = incoming.CallbackQuery.Message.Chat.ID
	msgId = incoming.CallbackQuery.Message.MessageID

	return
}

func (h tgHandler) getUserAndChatIDs(incoming tgbotapi.Update) (userId, chatId int64) {
	userId = incoming.SentFrom().ID
	chatId = incoming.Message.Chat.ID

	return
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
