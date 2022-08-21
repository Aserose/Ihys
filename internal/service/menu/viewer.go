package menu

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/service/webapi"
	"IhysBestowal/internal/service/webapi/tg/menu"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
	"strconv"
	"strings"
)

type enumContent func(song string, page int) []menu.Button
type contentWithControls func(enumContent) []menu.Button

type viewer struct {
	viewController
	viewAudio
}

func newViewer(cfg config.Keypads, md middleware, api webapi.WebApi) viewer {
	v := viewer{viewAudio: newViewAudio(cfg, md, api)}

	back := md.menu.LBtn(backTxt, backClb, func(p dto.Response) {
		md.menu.Build(v.newMsg(convert(p.MsgText).Song[0], p.ChatId), p, v.menuButtons(v.openContentWithControls)...)
	})

	v.viewController = newViewController(back, md)

	return v
}

func (v viewer) enumContent(src string, page int) []menu.Button {
	sml := v.md.items.get(src, page)
	b := make([]menu.Button, len(sml))

	for i, s := range sml {
		num := i
		b[i] = v.middleware.menu.LBtn(s.WithoutFrom(), strconv.Itoa(page+num), func(p dto.Response) {
			p.MsgId = 0
			v.openSong(p, v.md.get(p.MsgText, page)[num])
		})
	}

	return b
}

func (v viewer) openSong(p dto.Response, src datastruct.Song) {
	v.middleware.menu.Build(v.newMsg(src, p.ChatId), p, v.menuButtons(v.openContentWithControls)...)
}

func (v viewer) openContentWithControls(srcSong string, p dto.Response) {
	p.MsgText = srcSong
	v.build(false, v.enumContent, p)
}

func convert(msgTxt string) datastruct.Set {
	lSep, rSep := datastruct.Song{}.Separators()
	song := strings.Split(msgTxt, ` - `)

	if !strings.Contains(msgTxt, lSep) {
		return datastruct.Set{
			Song: []datastruct.Song{
				{
					Artist: song[0],
					Title:  strings.Split(song[1], dblIdt)[0],
				},
			},
		}
	}

	title := strings.Split(song[1], lSep)

	return datastruct.Set{
		From: strings.Replace(title[1], rSep, emp, 1),
		Song: []datastruct.Song{
			{
				Artist: song[0],
				Title:  title[0],
			},
		},
	}
}

func (v viewer) preload(p dto.Response) {
	for i := 0; i < 100; i++ {
		v.md.menu.Build(tgbotapi.MessageConfig{}, p, v.scroller[0][i](v.enumContent)...)
		v.md.menu.Build(tgbotapi.MessageConfig{}, p, v.scroller[1][i](v.enumContent)...)
		v.md.menu.Build(tgbotapi.MessageConfig{}, p, v.scroller[2][i](v.enumContent)...)
	}
}

func random(min, max int) int { return rand.Intn(max-min+1) + min }
