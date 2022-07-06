package menu

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/service/webapi"
	"IhysBestowal/internal/service/webapi/tg"
	"math/rand"
	"strconv"
	"strings"
)

type getEnumeratedContent func(sourceName string, page int) []tg.Button
type getContentWithControls func(getEnumeratedContent) []tg.Button

type viewer struct {
	viewController
	viewAudio
}

func newViewer(cfg config.Keypads, md middleware, api webapi.WebApiService) viewer {
	v := viewer{viewAudio: newViewItems(cfg, md, api)}

	backButton := md.tgBuilder.NewLineMenuButton(backText, backCallback, func(p dto.Response) {
		md.tgBuilder.MenuBuild(v.getSongMsgCfg(convert(p.MsgText).Items[0], p.ChatId), p, v.getSongMenuButtons(v.openContentListWithControls)...)
	})

	v.viewController = newViewController(backButton, md)

	return v
}

func (v viewer) getEnumeratedContent(sourceAudio string, page int) []tg.Button {
	items := v.middleware.items.get(sourceAudio, page)
	audioButtons := make([]tg.Button, len(items))

	for i, song := range items {
		num := i
		audioButtons[i] = v.middleware.tgBuilder.NewLineMenuButton(song.GetAudio(), strconv.Itoa(page+num), func(p dto.Response) {
			p.MsgId = 0
			v.openSongMenu(p, v.md.get(p.MsgText, page)[num])
		})
	}

	return audioButtons
}

func (v viewer) openSongMenu(p dto.Response, source datastruct.AudioItem) {
	v.middleware.tgBuilder.MenuBuild(v.getSongMsgCfg(source, p.ChatId), p, v.getSongMenuButtons(v.openContentListWithControls)...)
}

func (v viewer) openContentListWithControls(sourceSong string, p dto.Response) {
	p.MsgText = sourceSong
	v.buildMenu(false, v.getEnumeratedContent, p)
}

func convert(msgText string) datastruct.AudioItems {
	leftSepar, rightSepar := datastruct.AudioItem{}.GetSeparators()
	song := strings.Split(msgText, ` - `)

	if !strings.Contains(msgText, leftSepar) {
		return datastruct.AudioItems{
			Items: []datastruct.AudioItem{
				{
					Artist: song[0],
					Title:  strings.Split(song[1], doubleIndent)[0],
				},
			},
		}
	}

	s := strings.Split(song[1], leftSepar)

	return datastruct.AudioItems{
		From: strings.Replace(s[1], rightSepar, empty, 1),
		Items: []datastruct.AudioItem{
			{
				Artist: song[0],
				Title:  s[0],
			},
		},
	}
}

func getRandomNum(min, max int) int {
	return rand.Intn(max-min+1) + min
}
