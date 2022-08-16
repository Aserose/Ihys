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

	backButton := md.menu.NewLineMenuButton(backText, backCallback, func(p dto.Response) {
		md.menu.Build(v.msgCfg(convert(p.MsgText).Songs[0], p.ChatId), p, v.menuButtons(v.openContentListWithControls)...)
	})

	v.viewController = newViewController(backButton, md)

	return v
}

func (v viewer) enumContent(song string, page int) []menu.Button {
	it := v.md.items.get(song, page)
	audioButtons := make([]menu.Button, len(it))

	for i, song := range it {
		num := i
		audioButtons[i] = v.middleware.menu.NewLineMenuButton(song.WithoutFrom(), strconv.Itoa(page+num), func(p dto.Response) {
			p.MsgId = 0
			v.openSongMenu(p, v.md.get(p.MsgText, page)[num])
		})
	}

	return audioButtons
}

func (v viewer) openSongMenu(p dto.Response, src datastruct.Song) {
	v.middleware.menu.Build(v.msgCfg(src, p.ChatId), p, v.menuButtons(v.openContentListWithControls)...)
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
					Title:  strings.Split(song[1], doubleIndent)[0],
				},
			},
		}
	}

	s := strings.Split(song[1], leftSep)

	return datastruct.Songs{
		From: strings.Replace(s[1], rightSep, empty, 1),
		Songs: []datastruct.Song{
			{
				Artist: song[0],
				Title:  s[0],
			},
		},
	}
}

func (v viewer) init(p dto.Response) { v.setup(p, v.enumContent) }

func random(min, max int) int { return rand.Intn(max-min+1) + min }
