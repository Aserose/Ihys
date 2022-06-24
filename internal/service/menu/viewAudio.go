package menu

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/service/webapi"
	"IhysBestowal/internal/service/webapi/tg"
	"github.com/biter777/countries"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const (
	emojiMovieCamera  = " \xF0\x9F\x8E\xA5 "
	emojiLink         = " \xF0\x9F\x94\x97 "
	emojiHourglass    = " \xE2\x8C\x9B "
	emojiBlackNim     = " \xE2\x9C\x92 "
	emojiPageWithCurl = " \xF0\x9F\x93\x83 "

	separator      = ` | `
	msgYouTube     = separator + emojiMovieCamera + `[YouTube]`
	msgWebsite     = separator + emojiLink + `[Website]`
	msgLyrics      = separator + emojiPageWithCurl + `[Lyrics]`
	msgLoadingBase = emojiHourglass + `Un momento! It's uploading.`
)

var msgLoading = []string{}

func init() {
	msgLoading = append(msgLoading,
		newFormattedQuote(`Patience attracts happiness; it brings near that which is far.`, `Swahili Proverb`),
		newFormattedQuote(`Our patience will achieve more than our force.`, `Edmund Burke`),
		newFormattedQuote(`Learning patience can be a difficult experience, but once conquered, you will find life is easier.`, `Catherine Pulsifer`),
		newFormattedQuote(`Patience is the best remedy for every trouble.`, `Plautus`),
		newFormattedQuote(`Trees that are slow to grow bear the best fruit`, `Moliere`))
}

func newFormattedQuote(quote, author string) string {
	return msgLoadingBase + "\n\n" + "“" + quote + "“ \n - " + author + emojiBlackNim
}

type viewAudio struct {
	api webapi.WebApiService
	cfg config.Keypads
	middleware
}

func newViewItems(cfg config.Keypads, md middleware, api webapi.WebApiService) viewAudio {
	rand.Seed(time.Now().UnixNano())
	return viewAudio{
		api:        api,
		cfg:        cfg,
		middleware: md,
	}
}

func (vi viewAudio) getSongMsgCfg(song datastruct.AudioItem, chatId int64) tgbotapi.MessageConfig {
	songName := song.GetAudio()
	wg := &sync.WaitGroup{}

	resp := tgbotapi.NewMessage(chatId, " ")
	resp.ParseMode = `markdown`
	songName = "\n" + "[" + songName + "]" + "(" + song.Url + ")"
	resp.Text = songName + "\n\n"

	var (
		info      string
		ytURL     string
		website   string
		lyricsURL string
	)

	wg.Add(4)
	go func() {
		defer wg.Done()
		if songInfo := vi.api.GetSongInfo(song); songInfo.ReleaseDate != empty {
			var flg string
			if country := countries.ByName(songInfo.Country); country.String() != countries.UnknownMsg {
				flg = country.Emoji()
			}

			info =
				`Label: ` + songInfo.Label + ` < ` + songInfo.Country + `  ` + flg + ` > ` +
					"\n" + `Release: ` + songInfo.ReleaseDate +
					"\n" + `Genre: ` + songInfo.GetGenresString() +
					"\n\n"
		}
	}()
	go func() {
		defer wg.Done()
		if yTurl := vi.api.IYouTube.GetYTUrl(songName); yTurl != empty {
			ytURL = msgYouTube + vi.newFormattedURL(yTurl)
		}
	}()
	go func() {
		defer wg.Done()
		if web := vi.api.IDiscogs.GetWebsiteArtist(song.Artist); web != empty {
			website = msgWebsite + vi.newFormattedURL(web)
		}
	}()
	go func() {
		defer wg.Done()
		if lyrics := vi.api.IGenius.GetLyricsURL(song); lyrics != empty {
			lyricsURL = msgLyrics + vi.newFormattedURL(lyrics)
		}
	}()
	wg.Wait()

	resp.Text += info + ytURL + website + lyricsURL

	return resp
}

func (v viewAudio) newFormattedURL(url string) string {
	return `(` + url + `)`
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
				vi.api.Send(tgbotapi.NewEditMessageText(p.ChatId, p.MsgId, msgLoading[getRandomNum(0, len(msgLoading)-1)]))
				openMenu(vi.middleware.getAllSimilar(source), p)
			}),

		vi.tgBuilder.NewMenuButton(
			vi.cfg.SongMenu.Best.Text,
			vi.cfg.SongMenu.Best.CallbackData,
			func(p dto.Response) {
				source := convert(p.MsgText)
				vi.api.Send(tgbotapi.NewEditMessageText(p.ChatId, p.MsgId, msgLoading[getRandomNum(0, len(msgLoading)-1)]))
				openMenu(vi.middleware.getLastFMBest(source), p)
			}),
	}
}

func convert(msgText string) datastruct.AudioItems {
	song := strings.Split(msgText, ` - `)

	if !strings.Contains(msgText, ` «(`) {
		return datastruct.AudioItems{
			Items: []datastruct.AudioItem{
				{
					Artist: song[0],
					Title:  strings.Split(song[1], "\n\n")[0],
				},
			},
		}
	}

	s := strings.Split(song[1], ` «(`)

	return datastruct.AudioItems{
		From: strings.Replace(s[1], `)»`, ``, 1),
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
