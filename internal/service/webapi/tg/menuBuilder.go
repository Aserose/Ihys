package tg

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/dto"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TGMenu interface {
	MenuBuild(msg tgbotapi.MessageConfig, msgId int, execCmd dto.ExecCmd, menus ...Button)
	NewSubMenu(text, data string, menus ...Button) Button
	NewMenuButton(text, data string, onTapped func(chatId int64, msgId int)) Button
	NewLineSubMenu(text, data string, menus ...Button) Button
	NewLineSubMenuTap(text, data string, onTapped func(chatId int64, msgId int), menus ...Button) Button
	NewLineMenuButton(text, data string, onTapped func(chatId int64, msgId int)) Button
	IButton
}

type IButton interface {
	GetButtonText() string
	GetButtonData() string
}

type Builder struct {
	Button
	Api
	cfg config.Buttons
}

type Button struct {
	text         string
	callbackData string
	onTapped     func(chatId int64, msgId int)
	newline      bool
	menus        []Button
}

func newMenuBuilder(api Api, cfg config.Buttons) TGMenu {
	return Builder{
		Api: api,
		cfg: cfg,
	}
}

func (m Builder) MenuBuild(msg tgbotapi.MessageConfig, msgId int, execCmd dto.ExecCmd, menus ...Button) {
	var row [][]tgbotapi.InlineKeyboardButton
	mainMenu := m.cfg.MainMenu
	searchMenu := m.cfg.SearchMenu
	//TODO redo
	tempKludge := map[string]bool{
		mainMenu.CallbackData:   false,
		searchMenu.CallbackData: false,
	}

	for _, mmm := range menus {
		ms := mmm

		if ms.callbackData == mainMenu.CallbackData {
			tempKludge[mainMenu.CallbackData] = true
		}
		if ms.callbackData == searchMenu.CallbackData {
			tempKludge[searchMenu.CallbackData] = true
		}

		if ms.menus != nil || len(ms.menus) != 0 {
			ms.onTapped = func(chatId int64, msgId int) {
				m.MenuBuild(msg, msgId, execCmd, m.nextSubmenu(msg, execCmd, ms.menus, menus...)...)
			}
		}

		if ms.newline {
			row = append(row, []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(ms.text, ms.callbackData)})
		} else {
			if len(row) == 0 {
				row = append(row, []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(ms.text, ms.callbackData)})
			} else {
				row[len(row)-1] = append(row[len(row)-1], tgbotapi.NewInlineKeyboardButtonData(ms.text, ms.callbackData))
			}
		}
		execCmd[ms.callbackData] = ms.onTapped
	}

	if msgId != 0 {
		edt := tgbotapi.NewEditMessageTextAndMarkup(msg.ChatID, msgId, msg.Text, tgbotapi.NewInlineKeyboardMarkup(row...))
		edt.ParseMode = msg.ParseMode
		m.Api.Send(edt)
	} else {
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(row...)
		respId := m.Api.Send(msg).MessageID
		//TODO redo
		if tempKludge[mainMenu.CallbackData] {
			execCmd[mainMenu.Delete] = func(chatId int64, msgId int) {
				m.Api.Send(tgbotapi.NewDeleteMessage(chatId, respId))
			}
		}
		if tempKludge[searchMenu.CallbackData] {
			execCmd[searchMenu.Delete] = func(chatId int64, msgId int) {
				m.Api.Send(tgbotapi.NewDeleteMessage(chatId, respId))
			}
		}
	}
}

func (m Builder) nextSubmenu(resp tgbotapi.MessageConfig, execCmd dto.ExecCmd, submenus []Button, prevMenus ...Button) (res []Button) {
	for _, ms := range submenus {
		res = append(res, ms)
	}

	//res = append(res, Button{
	//	text:         "back",
	//	callbackData: fmt.Sprint("submenuBack", resp.ChatID, resp.Text),
	//	onTapped: func(chatId int64, msgId int) {
	//		m.MenuBuild(resp, msgId, execCmd, prevMenus...)
	//	}},
	//)

	execCmd[res[len(res)-1].callbackData] = res[len(res)-1].onTapped

	return
}

func (m Builder) NewLineSubMenuTap(text, callbackData string, onTapped func(chatId int64, msgId int), menus ...Button) Button {
	subMenu := Button{
		text:         text,
		callbackData: callbackData,
		onTapped:     onTapped,
		newline:      true,
	}

	for _, mm := range menus {
		subMenu.menus = append(subMenu.menus, mm)
	}

	return subMenu
}

func (m Builder) NewSubMenu(name, data string, menus ...Button) Button {
	subMenu := Button{
		text:         name,
		callbackData: data,
		onTapped:     nil,
	}

	for _, mm := range menus {
		subMenu.menus = append(subMenu.menus, mm)
	}

	return subMenu
}

func (m Builder) NewMenuButton(text, callbackData string, onTapped func(chatId int64, msgId int)) Button {
	return Button{
		text:         text,
		callbackData: callbackData,
		onTapped:     onTapped,
	}
}

func (m Builder) NewLineSubMenu(text, callbackData string, menus ...Button) Button {
	subMenu := Button{
		text:         text,
		callbackData: callbackData,
		newline:      true,
	}

	for _, mm := range menus {
		subMenu.menus = append(subMenu.menus, mm)
	}

	return subMenu
}

func (m Button) NewLineMenuButton(text, callbackData string, onTapped func(chatId int64, msgId int)) Button {
	return Button{
		text:         text,
		callbackData: callbackData,
		onTapped:     onTapped,
		newline:      true,
	}
}

func (m Button) GetButtonText() string {
	return m.text
}

func (m Button) GetButtonData() string {
	return m.callbackData
}
