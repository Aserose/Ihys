package menu

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/service/webapi"
	"IhysBestowal/internal/service/webapi/tg"
	"strconv"
)

type getEnumeratedContent func(sourceName string, page int) []tg.Button
type getContentWithControls func(getEnumeratedContent) []tg.Button

type viewer struct {
	viewController
	viewItems
}

func newViewer(cfg config.Menu, md middleware, api webapi.WebApiService) viewer {
	v := viewer{viewItems: newViewItems(cfg, md, api)}

	backButton := md.tgBuilder.NewLineMenuButton(trackMenu, back, func(p dto.Response) {
		md.tgBuilder.MenuBuild(v.getSongMsgCfg(convert(p.MsgText).Items[0], p.ChatId), p, v.getSongMenuButtons(v.openMenuWithControls)...)
	})

	v.viewController = newViewController(backButton, md)

	return v
}

func (v viewer) getEnumeratedContent(sourceName string, page int) []tg.Button {
	items := v.middleware.items.get(sourceName, page)
	audioButtons := make([]tg.Button, len(items))

	for i, song := range items {
		num := i
		audioButtons[i] = v.middleware.tgBuilder.NewLineMenuButton(song.Artist+` - `+song.Title, strconv.Itoa(page+num), func(p dto.Response) {
			p.MsgId = 0
			v.middleware.tgBuilder.MenuBuild(v.getSongMsgCfg(v.md.get(p.MsgText, page)[num], p.ChatId), p, v.getSongMenuButtons(v.openMenuWithControls)...)
		})
	}
	return audioButtons
}

func (v viewer) openMenuWithControls(sourceSong string, p dto.Response) {
	p.MsgText = sourceSong
	v.buildMenu(false, v.getEnumeratedContent, p)
}
