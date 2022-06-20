package menu

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/service/webapi"
	"IhysBestowal/internal/service/webapi/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type viewAudio struct {
	api webapi.WebApiService
	cfg config.Keypads
	middleware
}

func newViewItems(cfg config.Keypads, md middleware, api webapi.WebApiService) viewAudio {
	return viewAudio{
		api:        api,
		cfg:        cfg,
		middleware: md,
	}
}

func (vi viewAudio) getSongMsgCfg(song datastruct.AudioItem, chatId int64) tgbotapi.MessageConfig {
	songName := song.GetAudio()

	resp := tgbotapi.NewMessage(chatId, " ")
	resp.ParseMode = `markdown`
	YTurl := vi.api.IYouTube.GetYTUrl(songName)
	songName = "\n" + "[" + songName + "]" + "(" + song.Url + ")"
	if YTurl != " " {
		resp.Text = songName + "\n\n" + "\xF0\x9F\x8E\xA5 " + "[YouTube](" + YTurl + ")" + "\n\n"
	} else {
		resp.Text = songName
	}

	return resp
}

func (vi viewAudio) getSongMenuButtons(openMenu func(sourceName string, p dto.Response)) []tg.Button {
	return []tg.Button{

		vi.tgBuilder.NewMenuButton(
			vi.cfg.SongMenu.Delete.Text,
			vi.cfg.SongMenu.Delete.CallbackData,
			func(p dto.Response) {
				vi.api.Send(tgbotapi.NewDeleteMessage(p.ChatId, p.MsgId))
			}),

		vi.tgBuilder.NewMenuButton(
			vi.cfg.SongMenu.Similar.Text,
			vi.cfg.SongMenu.Similar.CallbackData,
			func(p dto.Response) {
				source := convert(p.MsgText)
				vi.api.Send(tgbotapi.NewEditMessageText(p.ChatId, p.MsgId, source.Items[0].GetAudio()))
				openMenu(vi.middleware.getAllSimilar(source), p)
			}),

		vi.tgBuilder.NewMenuButton(
			vi.cfg.SongMenu.Best.Text,
			vi.cfg.SongMenu.Best.CallbackData,
			func(p dto.Response) {
				source := convert(p.MsgText)
				vi.api.Send(tgbotapi.NewEditMessageText(p.ChatId, p.MsgId, source.Items[0].GetAudio()))
				openMenu(vi.middleware.getLastFMBest(source), p)
			}),
	}
}

func convert(msgText string) datastruct.AudioItems {
	song := strings.Split(msgText, ` - `)

	if !strings.Contains(msgText, ` (`) {
		return datastruct.AudioItems{
			Items: []datastruct.AudioItem{
				{
					Artist: song[0],
					Title:  strings.Split(song[1], "\n")[0],
				},
			},
		}
	}

	s := strings.Split(song[1], ` (`)

	return datastruct.AudioItems{
		From: strings.Replace(s[1], `)`, ``, 1),
		Items: []datastruct.AudioItem{
			{
				Artist: song[0],
				Title:  s[0],
			},
		},
	}
}
