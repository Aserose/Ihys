package menu

import (
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/service/webapi/tg/menu"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

type viewController struct {
	scroller [3][100]contentWithControls
	back     menu.Button
	lineSize int
	md       middleware
}

func newViewController(back menu.Button, md middleware) viewController {
	v := viewController{
		lineSize: 5,
		back:     back,
		md:       md,
	}

	for i := 0; i < 100; i++ {
		capt := i
		page := strconv.Itoa(capt + 1)

		vLeftClb := page + spc + vLeft
		vRightClb := page + spc + vRight
		vSelectClb := page + spc + vSelect

		v.scroller[0][i] = func(c enumContent) []menu.Button {
			return []menu.Button{
				v.md.menu.LBtn(leftArrow, vLeftClb, v.left(c)),
				v.md.menu.Btn(page, vSelectClb, v.openSelection(c)),
				v.md.menu.Btn(rightArrow, vRightClb, v.right(c)),
				v.back,
			}
		}

		v.scroller[1][i] = func(c enumContent) []menu.Button {
			return []menu.Button{
				v.md.menu.LBtn(page, vSelectClb, v.openSelection(c)),
				v.md.menu.Btn(rightArrow, vRightClb, v.right(c)),
				v.back,
			}
		}

		v.scroller[2][i] = func(c enumContent) []menu.Button {
			return []menu.Button{
				v.md.menu.LBtn(leftArrow, vLeftClb, v.left(c)),
				v.md.menu.Btn(page, vSelectClb, v.openSelection(c)),
				v.back,
			}
		}

	}

	return v
}

func (v viewController) isFirst(page int) bool {
	return page <= 0
}

func (v viewController) isLast(page int, src string) bool {
	return page >= v.md.pageCount(src)
}

func (v viewController) pageControls(page int, msgText string, c enumContent) []menu.Button {
	switch v.isFirst(page) {
	case true:
		switch v.isLast(page, msgText) {
		case true:
			return []menu.Button{v.back}
		case false:
			return v.scroller[1][page](c)
		}
	case false:
		switch v.isLast(page, msgText) {
		case true:
			return v.scroller[2][page](c)
		case false:

		}
	}
	return v.scroller[0][page](c)
}

func (v viewController) build(isLeft bool, c enumContent, p dto.Response) {
	msg := tgbotapi.MessageConfig{BaseChat: tgbotapi.BaseChat{ChatID: p.ChatId}, Text: p.MsgText}

	numPage, _ := strconv.Atoi(strings.Split(p.CallbackData, spc)[0])
	if isLeft {
		numPage -= 2
	}

	v.md.menu.Build(msg, p, append(c(p.MsgText, numPage), v.pageControls(numPage, p.MsgText, c)...)...)
}

func (v viewController) openSelection(c enumContent) dto.OnTappedFunc {
	return func(p dto.Response) {
		msg := tgbotapi.MessageConfig{BaseChat: tgbotapi.BaseChat{ChatID: p.ChatId}, Text: p.MsgText}
		v.md.menu.Build(msg, p, v.selection(c, p.MsgText)...)
	}
}

func (v viewController) selection(c enumContent, songMsgTxt string) []menu.Button {
	count := v.md.pageCount(songMsgTxt)
	slc := make([]menu.Button, count+1)
	isEndLine := func(elemNum int) bool { return elemNum%v.lineSize == 0 }

	for i := 0; i <= count; i++ {
		page := i + 1
		tap := func(p dto.Response) {
			v.build(false, c, p)
		}

		if isEndLine(i) {
			slc[i] = v.md.menu.LBtn(strconv.Itoa(page), strconv.Itoa(i)+spc+pageNum, tap)
		} else {
			slc[i] = v.md.menu.Btn(strconv.Itoa(page), strconv.Itoa(i)+spc+pageNum, tap)
		}
	}

	return slc
}

func (v viewController) right(c enumContent) dto.OnTappedFunc {
	return func(p dto.Response) {
		v.build(false, c, p)
	}
}

func (v viewController) left(c enumContent) dto.OnTappedFunc {
	return func(p dto.Response) {
		v.build(true, c, p)
	}
}
