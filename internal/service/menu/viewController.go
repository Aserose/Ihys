package menu

import (
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/service/webapi/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

const (
	rightArrow    = "->"
	leftArrow     = "<-"
	viewLeft      = "viewLeft"
	viewRight     = "viewRight"
	viewSelection = "viewSelection"
	pageNumber    = "pageNumber"
)

type viewController struct {
	scroller [][]getContentWithControls
	back     tg.Button
	lineSize int
	md       middleware
}

func newViewController(back tg.Button, md middleware) viewController {
	v := viewController{
		lineSize: 5,
		back:     back,
		md:       md,
	}
	v.scroller = make([][]getContentWithControls, 3)
	for i := 0; i < 3; i++ {
		v.scroller[i] = make([]getContentWithControls, 100)
	}

	for i := 0; i < 100; i++ {
		b := i
		numberPage := strconv.Itoa(b + 1)

		viewLeftCallback := numberPage + ` ` + viewLeft
		viewRightCallback := numberPage + ` ` + viewRight
		viewSelectionCallback := numberPage + ` ` + viewSelection

		v.scroller[0][i] = func(c getEnumeratedContent) []tg.Button {
			return []tg.Button{
				v.md.tgBuilder.NewLineMenuButton(leftArrow, viewLeftCallback, v.swapLeft(c)),
				v.md.tgBuilder.NewMenuButton(numberPage, viewSelectionCallback, v.openSelection(c)),
				v.md.tgBuilder.NewMenuButton(rightArrow, viewRightCallback, v.swapRight(c)),
				v.back,
			}
		}

		v.scroller[1][i] = func(c getEnumeratedContent) []tg.Button {
			return []tg.Button{
				v.md.tgBuilder.NewLineMenuButton(numberPage, viewSelectionCallback, v.openSelection(c)),
				v.md.tgBuilder.NewMenuButton(rightArrow, viewRightCallback, v.swapRight(c)),
				v.back,
			}
		}

		v.scroller[2][i] = func(c getEnumeratedContent) []tg.Button {
			return []tg.Button{
				v.md.tgBuilder.NewLineMenuButton(leftArrow, viewLeftCallback, v.swapLeft(c)),
				v.md.tgBuilder.NewMenuButton(numberPage, viewSelectionCallback, v.openSelection(c)),
				v.back,
			}
		}

	}

	return v
}

func (v viewController) isFirstPage(page int) bool {
	return page <= 0
}

func (v viewController) isLastPage(page int, sourceName string) bool {
	return page >= v.md.pageCount(sourceName)
}

func (v viewController) getPageControls(page int, msgText string, c getEnumeratedContent) []tg.Button {
	switch v.isFirstPage(page) {
	case true:
		switch v.isLastPage(page, msgText) {
		case true:
			return []tg.Button{v.back}
		case false:
			return v.scroller[1][page](c)
		}
	case false:
		switch v.isLastPage(page, msgText) {
		case true:
			return v.scroller[2][page](c)
		case false:

		}
	}
	return v.scroller[0][page](c)
}

func (v viewController) buildMenu(scrollBack bool, c getEnumeratedContent, p dto.Response) {
	msgCfg := tgbotapi.MessageConfig{BaseChat: tgbotapi.BaseChat{ChatID: p.ChatId}, Text: p.MsgText}

	page, _ := strconv.Atoi(strings.Split(p.CallbackData, ` `)[0])
	if scrollBack {
		page = page - 2
	}

	buttons := c(p.MsgText, page)

	buttons = append(buttons, v.getPageControls(page, p.MsgText, c)...)

	v.md.tgBuilder.MenuBuild(msgCfg, p, buttons...)
}

func (v viewController) openSelection(c getEnumeratedContent) dto.OnTappedFunc {
	return func(p dto.Response) {
		msgCfg := tgbotapi.MessageConfig{BaseChat: tgbotapi.BaseChat{ChatID: p.ChatId}, Text: p.MsgText}
		v.md.tgBuilder.MenuBuild(msgCfg, p, v.selection(c, p.MsgText)...)
	}
}

func (v viewController) selection(c getEnumeratedContent, songMsgTxt string) []tg.Button {
	pageAmount := v.md.pageCount(songMsgTxt)
	pageSelectionSubmenu := make([]tg.Button, pageAmount+1)
	isEndOfTheLine := func(elementNumber int) bool { return elementNumber%v.lineSize == 0 }

	for elementIndex := 0; elementIndex <= pageAmount; elementIndex++ {
		pageNum := elementIndex + 1
		selectButtonTapFunc := func(p dto.Response) {
			v.buildMenu(false, c, p)
		}

		if isEndOfTheLine(elementIndex) {
			pageSelectionSubmenu[elementIndex] = v.md.tgBuilder.NewLineMenuButton(strconv.Itoa(pageNum), strconv.Itoa(elementIndex)+` `+pageNumber, selectButtonTapFunc)
		} else {
			pageSelectionSubmenu[elementIndex] = v.md.tgBuilder.NewMenuButton(strconv.Itoa(pageNum), strconv.Itoa(elementIndex)+` `+pageNumber, selectButtonTapFunc)
		}
	}

	return pageSelectionSubmenu
}

func (v viewController) swapRight(c getEnumeratedContent) dto.OnTappedFunc {
	return func(p dto.Response) {
		v.buildMenu(false, c, p)
	}
}

func (v viewController) swapLeft(c getEnumeratedContent) dto.OnTappedFunc {
	return func(p dto.Response) {
		v.buildMenu(true, c, p)
	}
}
