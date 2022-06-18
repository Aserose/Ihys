package tg

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/dto"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TGMenu interface {
	MenuBuild(msg tgbotapi.MessageConfig, p dto.Response, menus ...Button)
	NewSubMenu(text, data string, menus ...Button) Button
	NewSubMenuTap(name, data string, tapFunc dto.OnTappedFunc, menus ...Button) Button
	NewMenuButton(text, data string, tapFunc dto.OnTappedFunc) Button
	NewLineSubMenu(text, data string, menus ...Button) Button
	NewLineSubMenuTap(text, data string, tapFunc dto.OnTappedFunc, menus ...Button) Button
	NewLineMenuButton(text, data string, tapFunc dto.OnTappedFunc) Button
	IButton
}

type IButton interface {
	GetButtonText() string
	GetButtonData() string
}

type Builder struct {
	Button
	Api
	cfg config.Menu
}

type Button struct {
	text         string
	callbackData string
	onTapped     dto.OnTappedFunc
	newline      bool
	menus        []Button
}

func newMenuBuilder(api Api, cfg config.Menu) TGMenu {
	return Builder{
		Api: api,
		cfg: cfg,
	}
}

func (m Builder) MenuBuild(msgCfg tgbotapi.MessageConfig, p dto.Response, menus ...Button) {
	var row [][]tgbotapi.InlineKeyboardButton

	for _, mmm := range menus {
		ms := mmm

		if ms.menus != nil || len(ms.menus) != 0 {
			ms.onTapped = func(p dto.Response) {
				m.MenuBuild(msgCfg, p, m.nextSubmenu(p.ExecCmd, ms.menus)...)
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
		p.ExecCmd[ms.callbackData] = ms.onTapped
	}

	if p.MsgId != 0 {
		edt := tgbotapi.NewEditMessageTextAndMarkup(msgCfg.ChatID, p.MsgId, msgCfg.Text, tgbotapi.NewInlineKeyboardMarkup(row...))
		edt.ParseMode = msgCfg.ParseMode
		m.Api.Send(edt)
	} else {
		msgCfg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(row...)
		m.Api.Send(msgCfg)
	}
}

func (m Builder) nextSubmenu(execCmd dto.ExecCmd, submenus []Button) (res []Button) {
	for _, ms := range submenus {
		res = append(res, ms)
	}

	execCmd[res[len(res)-1].callbackData] = res[len(res)-1].onTapped

	return
}

func (m Builder) NewLineSubMenuTap(text, callbackData string, tapFunc dto.OnTappedFunc, menus ...Button) Button {
	subMenu := Button{
		text:         text,
		callbackData: callbackData,
		onTapped:     tapFunc,
		newline:      true,
	}

	for _, mm := range menus {
		subMenu.menus = append(subMenu.menus, mm)
	}

	return subMenu
}

func (m Builder) NewSubMenuTap(name, data string, tapFunc dto.OnTappedFunc, menus ...Button) Button {
	subMenu := Button{
		text:         name,
		callbackData: data,
		onTapped:     tapFunc,
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

func (m Builder) NewMenuButton(text, callbackData string, tapFunc dto.OnTappedFunc) Button {
	return Button{
		text:         text,
		callbackData: callbackData,
		onTapped:     tapFunc,
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

func (m Button) NewLineMenuButton(text, callbackData string, tapFunc dto.OnTappedFunc) Button {
	return Button{
		text:         text,
		callbackData: callbackData,
		onTapped:     tapFunc,
		newline:      true,
	}
}

func (m Button) GetButtonText() string {
	return m.text
}

func (m Button) GetButtonData() string {
	return m.callbackData
}
