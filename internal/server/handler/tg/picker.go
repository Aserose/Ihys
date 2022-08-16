package tg

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

// picker is responsible for retrieving certain parameters from the update response
type picker struct{}

func (p picker) exeArgs(incoming tgbotapi.Update) (callbackData string, chatId int64, msgId int) {
	callbackData = incoming.CallbackQuery.Data
	chatId, msgId = incoming.CallbackQuery.Message.Chat.ID, incoming.CallbackQuery.Message.MessageID

	return
}

func (p picker) cmdArgs(incoming tgbotapi.Update) (userId, chatId int64, query string) {
	userId, chatId = p.userAndChatIDs(incoming)
	query = p.cmdContent(incoming.Message.Text, incoming.Message.Command())

	return
}

func (p picker) callbackChatAndMsgIDs(incoming tgbotapi.Update) (chatId int64, msgId int) {
	chatId = incoming.CallbackQuery.Message.Chat.ID
	msgId = incoming.CallbackQuery.Message.MessageID

	return
}

func (p picker) userAndChatIDs(incoming tgbotapi.Update) (userId, chatId int64) {
	userId = incoming.SentFrom().ID
	chatId = incoming.Message.Chat.ID

	return
}

func (p picker) cmdContent(rawMsgText, nameCmd string) string {
	result := strings.Split(rawMsgText, nameCmd+" ")
	if len(result) <= 1 {
		return ``
	}
	return strings.Split(rawMsgText, nameCmd+" ")[1]
}

func (p picker) userAndMsgIDs(incoming tgbotapi.Update) (userId int64, msgId int) {
	userId = incoming.SentFrom().ID
	msgId = incoming.Message.MessageID

	return
}

func newExe(service service.Service) dto.ExecCmd {
	exe := make(map[string]dto.OnTappedFunc)
	exe[dlt] = func(p dto.Response) { service.TG.Send(tgbotapi.NewDeleteMessage(p.ChatId, p.MsgId)) }

	leftSep, rightSep := datastruct.Song{}.Separators()
	service.Menu.Setup(dto.Response{ExecCmd: exe, MsgText: `rick astley - never gonna give you up` + leftSep + `all` + rightSep})

	return exe
}
