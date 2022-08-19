package menu

import "IhysBestowal/internal/dto"

type Button struct {
	txt      string
	callback string
	onTapped dto.OnTappedFunc
	newline  bool
	menus    []Button
}

func (m Button) NewLineMenuButton(text, callback string, tap dto.OnTappedFunc) Button {
	return Button{
		txt:      text,
		callback: callback,
		onTapped: tap,
		newline:  true,
	}
}

func (m Button) Text() string     { return m.txt }
func (m Button) Callback() string { return m.callback }