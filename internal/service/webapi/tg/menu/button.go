package menu

import "IhysBestowal/internal/dto"

type Button struct {
	txt      string
	callback string
	onTapped dto.OnTappedFunc
	newline  bool
	menus    []Button
}

func (m Button) NewLineMenuButton(text, callbackData string, tap dto.OnTappedFunc) Button {
	return Button{
		txt:      text,
		callback: callbackData,
		onTapped: tap,
		newline:  true,
	}
}

func (m Button) Text() string { return m.txt }
func (m Button) Data() string { return m.callback }
