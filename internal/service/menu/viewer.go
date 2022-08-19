package menu

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/service/webapi"
	"IhysBestowal/internal/service/webapi/tg/menu"
	"math/rand"
	"strconv"
	"strings"
)

type enumContent func(song string, page int) []menu.Button
type contentWithControls func(enumContent) []menu.Button

type viewer struct {
	viewController
	viewSong
}

func newViewer(cfg config.Keypads, md middleware, api webapi.WebApi) viewer {
	v := viewer{viewSong: newViewItems(cfg, md, api)}

	backButton := md.menu.NewLineMenuButton(backTxt, backClbck, func(p dto.Response) {
		md.menu.Build(v.msg(convert(p.MsgText).Songs[0], p.ChatId), p, v.menuButtons(v.openContentListWithControls)...)
	})

	v.viewController = newViewController(backButton, md)

	return v
}

func (v viewer) enumContent(src string, page int) []menu.Button {
	sml := v.md.items.get(src, page)
	b := make([]menu.Button, len(sml))

	for i, s := range sml {
		num := i
		b[i] = v.middleware.menu.NewLineMenuButton(s.WithoutFrom(), strconv.Itoa(page+num), func(p dto.Response) {
			p.MsgId = 0
			v.openSongMenu(p, v.md.get(p.MsgText, page)[num])
		})
	}

	return b
}

func (v viewer) openSongMenu(p dto.Response, src datastruct.Song) {
	v.middleware.menu.Build(v.msg(src, p.ChatId), p, v.menuButtons(v.openContentListWithControls)...)
}

func (v viewer) openContentListWithControls(srcSong string, p dto.Response) {
	p.MsgText = srcSong
	v.build(false, v.enumContent, p)
}

func convert(msgTxt string) datastruct.Songs {
	leftSep, rightSep := datastruct.Song{}.Separators()
	song := strings.Split(msgTxt, ` - `)

	if !strings.Contains(msgTxt, leftSep) {
		return datastruct.Songs{
			Songs: []datastruct.Song{
				{
					Artist: song[0],
					Title:  strings.Split(song[1], dblIdt)[0],
				},
			},
		}
	}

	title := strings.Split(song[1], leftSep)

	return datastruct.Songs{
		From: strings.Replace(title[1], rightSep, emp, 1),
		Songs: []datastruct.Song{
			{
				Artist: song[0],
				Title:  title[0],
			},
		},
	}
}

func (v viewer) init(p dto.Response) { v.setup(p, v.enumContent) }

func random(min, max int) int { return rand.Intn(max-min+1) + min }
