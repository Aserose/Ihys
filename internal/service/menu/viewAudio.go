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

type viewSong struct {
	api webapi.WebApi
	cfg config.Keypads
	middleware
}

func newViewItems(cfg config.Keypads, md middleware, api webapi.WebApi) viewSong {
	return viewSong{
		api:        api,
		cfg:        cfg,
		middleware: md,
	}
}

func (vs viewSong) msgCfg(src datastruct.Song, chatId int64) tgbotapi.MessageConfig {
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
		if songInfo := vs.api.SongInfo(src); songInfo.ReleaseDate != emp {
			info = formatInfo(songInfo)
		}
	}()
	go func() {
		defer wg.Done()
		if vidURL := vs.api.YouTube.VideoURL(title); vidURL != emp {
			videoURL = formatVideoURL(vidURL)
		}
	}()
	go func() {
		defer wg.Done()
		if web := vs.api.Discogs.SiteArtist(src.Artist); web != emp {
			website = formatWebsite(web)
		}
	}()
	go func() {
		defer wg.Done()
		if lyrURL := vs.api.Genius.LyricsURL(src); lyrURL != emp {
			lyricsURL = formatLyricsURL(lyrURL)
		}
	}()
	wg.Wait()

	resp.Text += info + videoURL + website + lyricsURL + dblIdt

	return resp
}

func (vs viewSong) menuButtons(openMenu func(src string, p dto.Response)) []menu.Button {
	return []menu.Button{

		vs.menu.NewMenuButton(
			vs.cfg.SongMenu.Delete.Text,
			vs.cfg.SongMenu.Delete.CallbackData,
			func(p dto.Response) {
				vs.api.TG.Send(tgbotapi.NewDeleteMessage(p.ChatId, p.MsgId))
			}),

		vs.menu.NewMenuButton(
			vs.cfg.SongMenu.Similar.Text,
			vs.cfg.SongMenu.Similar.CallbackData,
			func(p dto.Response) {
				source := convert(p.MsgText)
				source.From = vs.middleware.from().All()

				vs.api.TG.Send(tgbotapi.NewEditMessageText(p.ChatId, p.MsgId, msgLoading[random(0, len(msgLoading)-1)]))
				openMenu(vs.middleware.similar(source), p)
			}),

		vs.menu.NewMenuButton(
			vs.cfg.SongMenu.Best.Text,
			vs.cfg.SongMenu.Best.CallbackData,
			func(p dto.Response) {
				source := convert(p.MsgText)
				source.From = vs.middleware.from().Lfm().Top()

				vs.api.TG.Send(tgbotapi.NewEditMessageText(p.ChatId, p.MsgId, msgLoading[random(0, len(msgLoading)-1)]))
				openMenu(vs.middleware.similar(source), p)
			}),
	}
}
