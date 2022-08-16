package menu

import (
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/service/webapi/tg/menu"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

type viewController struct {
	scroller [3][]contentWithControls
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

	v.scroller = [3][]contentWithControls{}

	for i := 0; i < 3; i++ {
		v.scroller[i] = make([]contentWithControls, 100)
	}

	for i := 0; i < 100; i++ {
		b := i
		numPage := strconv.Itoa(b + 1)

		viewLeftCallback := numPage + space + viewLeft
		viewRightCallback := numPage + space + viewRight
		viewSelectionCallback := numPage + space + viewSelection

		v.scroller[0][i] = func(c enumContent) []menu.Button {
			return []menu.Button{
				v.md.menu.NewLineMenuButton(leftArrow, viewLeftCallback, v.left(c)),
				v.md.menu.NewMenuButton(numPage, viewSelectionCallback, v.openSelection(c)),
				v.md.menu.NewMenuButton(rightArrow, viewRightCallback, v.right(c)),
				v.back,
			}
		}

		v.scroller[1][i] = func(c enumContent) []menu.Button {
			return []menu.Button{
				v.md.menu.NewLineMenuButton(numPage, viewSelectionCallback, v.openSelection(c)),
				v.md.menu.NewMenuButton(rightArrow, viewRightCallback, v.right(c)),
				v.back,
			}
		}

		v.scroller[2][i] = func(c enumContent) []menu.Button {
			return []menu.Button{
				v.md.menu.NewLineMenuButton(leftArrow, viewLeftCallback, v.left(c)),
				v.md.menu.NewMenuButton(numPage, viewSelectionCallback, v.openSelection(c)),
				v.back,
			}
		}

	}

	return v
}

func (v viewer) setup(p dto.Response, c enumContent) {
	for i := 0; i < 100; i++ {
		v.md.menu.Build(tgbotapi.MessageConfig{}, p, v.scroller[0][i](c)...)
		v.md.menu.Build(tgbotapi.MessageConfig{}, p, v.scroller[1][i](c)...)
		v.md.menu.Build(tgbotapi.MessageConfig{}, p, v.scroller[2][i](c)...)
	}
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

func (v viewController) build(isBack bool, ec enumContent, p dto.Response) {
	msgCfg := tgbotapi.MessageConfig{BaseChat: tgbotapi.BaseChat{ChatID: p.ChatId}, Text: p.MsgText}

	page, _ := strconv.Atoi(strings.Split(p.CallbackData, space)[0])
	if isBack {
		page = page - 2
	}

	btns := ec(p.MsgText, page)

	btns = append(btns, v.pageControls(page, p.MsgText, ec)...)

	v.md.menu.Build(msgCfg, p, btns...)
}

func (v viewController) openSelection(c enumContent) dto.OnTappedFunc {
	return func(p dto.Response) {
		msgCfg := tgbotapi.MessageConfig{BaseChat: tgbotapi.BaseChat{ChatID: p.ChatId}, Text: p.MsgText}
		v.md.menu.Build(msgCfg, p, v.selection(c, p.MsgText)...)
	}
}

func (v viewController) selection(c enumContent, songMsgTxt string) []menu.Button {
	pageAmount := v.md.pageCount(songMsgTxt)
	pageSelectionSubmenu := make([]menu.Button, pageAmount+1)
	isEndOfTheLine := func(elementNumber int) bool { return elementNumber%v.lineSize == 0 }

	for i := 0; i <= pageAmount; i++ {
		pageNum := i + 1
		selectButtonTapFunc := func(p dto.Response) {
			v.build(false, c, p)
		}

		if isEndOfTheLine(i) {
			pageSelectionSubmenu[i] = v.md.menu.NewLineMenuButton(strconv.Itoa(pageNum), strconv.Itoa(i)+space+pageNumber, selectButtonTapFunc)
		} else {
			pageSelectionSubmenu[i] = v.md.menu.NewMenuButton(strconv.Itoa(pageNum), strconv.Itoa(i)+space+pageNumber, selectButtonTapFunc)
		}
	}

	return pageSelectionSubmenu
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
