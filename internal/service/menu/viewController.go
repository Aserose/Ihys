package menu

import (
	"IhysBestowal/internal/service/webapi/tg"
	"strconv"
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
	setCurrentPage func(num int)
	getCurrentPage func() int
	back           tg.Button
	pageAmount     int
	lineSize       int
}

func newViewController(pageAmount int, back tg.Button) viewController {
	currentPage := 1

	return viewController{
		setCurrentPage: func(num int) { currentPage = num },
		pageAmount:     pageAmount,
		lineSize:       5,
		back:           back,
		getCurrentPage: func() int {
			return currentPage
		},
	}
}

func (v viewController) SetCurrentPageNumber(num int) {
	v.setCurrentPage(num)
}

func (v viewController) GetCurrentPageNumber() int {
	return v.getCurrentPage()
}

func (v viewController) IsFirstPage() bool {
	return v.GetCurrentPageNumber()-1 <= 0
}

func (v viewController) IsLastPage() bool {
	return v.GetCurrentPageNumber() > v.pageAmount
}

func (v viewController) GetControlButtons(callbackData string, swapLeft, swapRight bool, setPageContent func(), builder tg.TGMenu) []tg.Button {
	pageSelectionSubmenu := v.selection(setPageContent, builder)
	toSwapRight := v.swapRight(setPageContent)
	toSwapLeft := v.swapLeft(setPageContent)

	viewLeft := callbackData + viewLeft
	viewRight := callbackData + viewRight
	viewSelection := callbackData + viewSelection

	numberPage := strconv.Itoa(v.GetCurrentPageNumber())

	switch swapRight {
	case true:
		switch swapLeft {
		case true:
			return []tg.Button{
				builder.NewLineMenuButton(leftArrow, viewLeft, toSwapLeft),
				builder.NewSubMenu(numberPage, viewSelection, pageSelectionSubmenu...),
				builder.NewMenuButton(rightArrow, viewRight, toSwapRight),
				v.back,
			}
		case false:
			return []tg.Button{
				builder.NewLineSubMenu(numberPage, viewSelection, pageSelectionSubmenu...),
				builder.NewMenuButton(rightArrow, viewRight, toSwapRight),
				v.back,
			}
		}
	case false:
		switch swapLeft {
		case true:
			return []tg.Button{
				builder.NewLineMenuButton(leftArrow, viewLeft, toSwapLeft),
				builder.NewSubMenu(numberPage, viewSelection, pageSelectionSubmenu...),
				v.back,
			}
		case false:

		}
	}

	return []tg.Button{v.back}
}

func (v viewController) selection(fillContent func(), builder tg.TGMenu) []tg.Button {
	var (
		pageSelectionSubmenu = make([]tg.Button, v.pageAmount+1)
		isEndOfTheLine       = func(elementNumber int) bool { return elementNumber%v.lineSize == 0 }
	)

	for elementIndex := 0; elementIndex <= v.pageAmount; elementIndex++ {
		pageNum := elementIndex + 1
		selectButtonTapFunc := func(chatId int64, interMsgId int) {
			v.SetCurrentPageNumber(pageNum)
			fillContent()
		}

		if isEndOfTheLine(elementIndex) {
			pageSelectionSubmenu[elementIndex] = builder.NewLineMenuButton(strconv.Itoa(pageNum), pageNumber+strconv.Itoa(pageNum), selectButtonTapFunc)
		} else {
			pageSelectionSubmenu[elementIndex] = builder.NewMenuButton(strconv.Itoa(pageNum), pageNumber+strconv.Itoa(pageNum), selectButtonTapFunc)
		}
	}

	return pageSelectionSubmenu
}

func (v viewController) swapRight(setPageContent func()) func(chatId int64, interMsgId int) {
	return func(chatId int64, interMsgId int) {
		v.SetCurrentPageNumber(v.GetCurrentPageNumber() + 1)
		setPageContent()
	}
}

func (v viewController) swapLeft(setPageContent func()) func(chatId int64, interMsgId int) {
	return func(chatId int64, interMsgId int) {
		v.SetCurrentPageNumber(v.GetCurrentPageNumber() - 1)
		setPageContent()
	}
}
