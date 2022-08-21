package menu

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/service/webapi"
	"IhysBestowal/internal/service/webapi/tg/menu"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sync"
)

type viewAudio struct {
	api webapi.WebApi
	cfg config.Keypads
	middleware
}

func newViewAudio(cfg config.Keypads, md middleware, api webapi.WebApi) viewAudio {
	return viewAudio{
		api:        api,
		cfg:        cfg,
		middleware: md,
	}
}

func (va viewAudio) newMsg(src datastruct.Song, chatId int64) tgbotapi.MessageConfig {
	wg := &sync.WaitGroup{}

	resp := tgbotapi.NewMessage(chatId, " ")
	resp.ParseMode = `markdown`
	title := formatSong(src)
	resp.Text = title + "\n\n"

	var (
		info      string
		videoURL  string
		website   string
		lyricsURL string
	)

	wg.Add(4)
	go func() {
		defer wg.Done()
		if s := va.api.SongInfo(src); s.ReleaseDate != emp {
			info = formatInfo(s)
		}
	}()
	go func() {
		defer wg.Done()
		if v := va.api.YouTube.VideoURL(title); v != emp {
			videoURL = formatVideoURL(v)
		}
	}()
	go func() {
		defer wg.Done()
		if w := va.api.Discogs.SiteArtist(src.Artist); w != emp {
			website = formatWebsite(w)
		}
	}()
	go func() {
		defer wg.Done()
		if l := va.api.Genius.LyricsURL(src); l != emp {
			lyricsURL = formatLyricsURL(l)
		}
	}()
	wg.Wait()

	resp.Text += build(info, videoURL, website, lyricsURL, dblIdt)

	return resp
}

func (va viewAudio) menuButtons(content func(src string, p dto.Response)) []menu.Button {
	return []menu.Button{

		va.menu.Btn(
			va.cfg.SongMenu.Delete.Text,
			va.cfg.SongMenu.Delete.Callback,
			func(p dto.Response) {
				va.api.TG.Send(tgbotapi.NewDeleteMessage(p.ChatId, p.MsgId))
			}),

		va.menu.Btn(
			va.cfg.SongMenu.Similar.Text,
			va.cfg.SongMenu.Similar.Callback,
			func(p dto.Response) {
				src := convert(p.MsgText)
				src.From = va.middleware.from().All()

				va.api.TG.Send(tgbotapi.NewEditMessageText(p.ChatId, p.MsgId, msgLoading[random(0, len(msgLoading)-1)]))
				content(va.middleware.similar(src), p)
			}),

		va.menu.Btn(
			va.cfg.SongMenu.Best.Text,
			va.cfg.SongMenu.Best.Callback,
			func(p dto.Response) {
				src := convert(p.MsgText)
				src.From = va.middleware.from().LastFm().Top()

				va.api.TG.Send(tgbotapi.NewEditMessageText(p.ChatId, p.MsgId, msgLoading[random(0, len(msgLoading)-1)]))
				content(va.middleware.similar(src), p)
			}),
	}
}
