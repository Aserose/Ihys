package menu

import (
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/service/webapi/tg/api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Builder struct {
	Button
	api.Api
}

func New(api api.Api) Builder {
	return Builder{
		Api: api,
	}
}

func (m Builder) Build(msgCfg tgbotapi.MessageConfig, p dto.Response, menus ...Button) {
	var row [][]tgbotapi.InlineKeyboardButton

	for _, mn := range menus {
		ms := mn

		if ms.menus != nil || len(ms.menus) != 0 {
			ms.onTapped = func(p dto.Response) {
				m.Build(msgCfg, p, m.nextSubmenu(p.ExecCmd, ms.menus)...)
			}
		}

		if ms.newline {
			row = append(row, []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(ms.txt, ms.callback)})
		} else {
			if len(row) == 0 {
				row = append(row, []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(ms.txt, ms.callback)})
			} else {
				row[len(row)-1] = append(row[len(row)-1], tgbotapi.NewInlineKeyboardButtonData(ms.txt, ms.callback))
			}
		}

		p.ExecCmd[ms.callback] = ms.onTapped
	}

	if p.TGUser == (dto.TGUser{}) {
		return
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

func (m Builder) nextSubmenu(exec dto.ExecCmd, submenus []Button) (res []Button) {
	defer func() {
		exec[res[len(res)-1].callback] = res[len(res)-1].onTapped
	}()

	return submenus
}

func (m Builder) NewLineSubMenuTap(txt, callback string, tap dto.OnTappedFunc, menus ...Button) Button {
	return Button{
		txt:      txt,
		callback: callback,
		onTapped: tap,
		newline:  true,
		menus:    menus,
	}
}

func (m Builder) NewSubMenuTap(txt, callback string, tap dto.OnTappedFunc, menus ...Button) Button {
	return Button{
		txt:      txt,
		callback: callback,
		onTapped: tap,
		menus:    menus,
	}
}

func (m Builder) NewSubMenu(txt, callback string, menus ...Button) Button {
	return Button{
		txt:      txt,
		callback: callback,
		onTapped: nil,
		menus:    menus,
	}
}

func (m Builder) NewMenuButton(txt, callback string, tap dto.OnTappedFunc) Button {
	return Button{
		txt:      txt,
		callback: callback,
		onTapped: tap,
	}
}

func (m Builder) NewLineSubMenu(txt, callback string, menus ...Button) Button {
	return Button{
		txt:      txt,
		callback: callback,
		newline:  true,
		menus:    menus,
	}
}
