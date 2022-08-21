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

func (m Builder) Build(msgCfg tgbotapi.MessageConfig, p dto.Response, btn ...Button) {
	var row [][]tgbotapi.InlineKeyboardButton

	for _, bt := range btn {
		b := bt

		if b.btn != nil || len(b.btn) != 0 {
			b.onTapped = func(p dto.Response) {
				m.Build(msgCfg, p, m.nextSub(p.ExecCmd, b.btn)...)
			}
		}

		if b.newline {
			row = append(row, []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(b.txt, b.callback)})
		} else {
			if len(row) == 0 {
				row = append(row, []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(b.txt, b.callback)})
			} else {
				row[len(row)-1] = append(row[len(row)-1], tgbotapi.NewInlineKeyboardButtonData(b.txt, b.callback))
			}
		}

		p.ExecCmd[b.callback] = b.onTapped
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

func (m Builder) nextSub(exec dto.ExecCmd, btn []Button) (res []Button) {
	defer func() {
		exec[res[len(res)-1].callback] = res[len(res)-1].onTapped
	}()

	return btn
}

func (m Builder) LSubTap(txt, callback string, tap dto.OnTappedFunc, btn ...Button) Button {
	return Button{
		txt:      txt,
		callback: callback,
		onTapped: tap,
		newline:  true,
		btn:      btn,
	}
}

func (m Builder) SubTap(txt, callback string, tap dto.OnTappedFunc, btn ...Button) Button {
	return Button{
		txt:      txt,
		callback: callback,
		onTapped: tap,
		btn:      btn,
	}
}

func (m Builder) Sub(txt, callback string, btn ...Button) Button {
	return Button{
		txt:      txt,
		callback: callback,
		onTapped: nil,
		btn:      btn,
	}
}

func (m Builder) LBtn(txt, callback string, tap dto.OnTappedFunc) Button {
	return Button{
		txt:      txt,
		callback: callback,
		onTapped: tap,
		newline:  true,
	}
}

func (m Builder) Btn(txt, callback string, tap dto.OnTappedFunc) Button {
	return Button{
		txt:      txt,
		callback: callback,
		onTapped: tap,
	}
}

func (m Builder) LSub(txt, callback string, btn ...Button) Button {
	return Button{
		txt:      txt,
		callback: callback,
		newline:  true,
		btn:      btn,
	}
}
