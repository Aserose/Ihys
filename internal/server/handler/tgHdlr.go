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
	exe     map[string]func(chatId int64, msgId int)
	p       picker
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
		h.execute(incoming)

	default:
		h.deleteMsg(incoming)
		if incoming.Message.IsCommand() {
			switch incoming.Message.Command() {

			case "mainMenu":
				h.openMainMenu(incoming)

			case "authVk":
				h.authVk(incoming)

			case "start":
				// TODO

			case "search":
				h.search(incoming)

			}
		}
	}
}

func (h tgHandler) execute(incoming tgbotapi.Update) {
	callbackData, chatId, msgId := h.p.getExeArgs(incoming)
	if h.exe[callbackData] == nil {

		//TODO

	} else {
		h.exe[callbackData](chatId, msgId)
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
		ChatID: chatId,
	}, key)

	h.openMainMenu(incoming)
}

func (h tgHandler) openMainMenu(incoming tgbotapi.Update) {
	userId, chatId := h.p.getUserAndChatIDs(incoming)

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

func (h tgHandler) deleteMsg(incoming tgbotapi.Update) {
	chatId, msgId := h.p.getUserAndMsgIDs(incoming)
	h.exe["delete"](chatId, msgId)
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
