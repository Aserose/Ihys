package webapi

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/repository"
	"IhysBestowal/internal/service/auth"
	"IhysBestowal/internal/service/webapi/discogs"
	"IhysBestowal/internal/service/webapi/genius"
	"IhysBestowal/internal/service/webapi/gnoosic"
	"IhysBestowal/internal/service/webapi/lastFm"
	"IhysBestowal/internal/service/webapi/soundcloud"
	tg "IhysBestowal/internal/service/webapi/tg"
	"IhysBestowal/internal/service/webapi/tg/menu"
	"IhysBestowal/internal/service/webapi/vk"
	"IhysBestowal/internal/service/webapi/yaMusic"
	"IhysBestowal/internal/service/webapi/youTube"
	"IhysBestowal/pkg/customLogger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sort"
	"strings"
	"sync"
)

type Genius interface {
	LyricsURL(src datastruct.Song) string
}

type Gnoosic interface {
	RandomArtist() string
}

type Discogs interface {
	SiteLabel(src string) string
	SiteArtist(src string) string
	SongInfo(src datastruct.Song) datastruct.SongInfo
}

type Button interface {
	Text() string
	Data() string
}

type Menu interface {
	Build(msg tgbotapi.MessageConfig, p dto.Response, menus ...menu.Button)
	NewSubMenu(txt, callback string, menus ...menu.Button) menu.Button
	NewSubMenuTap(txt, callback string, tap dto.OnTappedFunc, menus ...menu.Button) menu.Button
	NewMenuButton(txt, callback string, tap dto.OnTappedFunc) menu.Button
	NewLineSubMenu(txt, callback string, menus ...menu.Button) menu.Button
	NewLineSubMenuTap(txt, callback string, tap dto.OnTappedFunc, menus ...menu.Button) menu.Button
	NewLineMenuButton(txt, callback string, tap dto.OnTappedFunc) menu.Button
	Button
}

type TG interface {
	Menu
	Send(c tgbotapi.Chattable) tgbotapi.Message
}

type LastFM interface {
	Auth(uid int64)
	Find(query string) datastruct.Song
	Similar(uid int64, src datastruct.Songs, opts ...lastFm.Set) datastruct.Songs
	Top(artists []string, max int) datastruct.Songs
}

type Soundcloud interface {
	Similar(src datastruct.Songs, opts ...soundcloud.Set) datastruct.Songs
	Close()
}

type VK interface {
	Recommendations(user dto.TGUser, offset int) (datastruct.Songs, error)
	RecommendationsCustom(user dto.TGUser) (datastruct.Songs, error)
	UserPlaylists(user dto.TGUser) (datastruct.Playlists, error)
	PlaylistSongs(user dto.TGUser, playlistId, ownerId int) (datastruct.Songs, error)
	VKAuth
}

type VKAuth interface {
	Auth(user dto.TGUser, serviceKey string) error
	AuthURL() string
	IsAuthorized(user dto.TGUser) bool
	IsValid(token string) bool
}

type YaMusic interface {
	Similar(src datastruct.Songs, opts ...yaMusic.Set) datastruct.Songs
	Find(query string) datastruct.Song
}

type YouTube interface {
	VideoURL(query string) string
}

type WebApi struct {
	VK
	TG
	YouTube
	LastFM
	YaMusic
	Soundcloud
	Discogs
	Gnoosic
	Genius
	From
}

func New(log customLogger.Logger, cfg config.Service, repo repository.Repository, authService auth.Auth) WebApi {
	return WebApi{
		TG:         tg.New(log, cfg),
		VK:         vk.New(log, cfg.Vk, authService.Vk()),
		YouTube:    youTube.New(log),
		LastFM:     lastFm.New(log, cfg.LastFM, repo),
		YaMusic:    yaMusic.New(log),
		Soundcloud: soundcloud.New(log),
		Discogs:    discogs.New(log, cfg.Discogs),
		Gnoosic:    gnoosic.New(),
		Genius:     genius.New(log, cfg.Genius),
		From:       newFrom(),
	}
}

func (s WebApi) Random() datastruct.Song {
	if item := s.Top(s.Gnoosic.RandomArtist()).Songs[0]; item != (datastruct.Song{}) {
		return item
	}

	return datastruct.Song{}
}

func (s WebApi) Find(query string) datastruct.Song {
	if response := s.LastFM.Find(query); response.Title != `` {
		return response
	}
	return s.YaMusic.Find(query)
}

func (s WebApi) Similar(src datastruct.Songs, opt Opt) datastruct.Songs {
	wg := &sync.WaitGroup{}
	res := []datastruct.Song{}
	ch := make(chan []datastruct.Song)
	cls := make(chan struct{})

	go func() {
		for {
			select {
			case i := <-ch:
				res = append(res, i...)
			case <-cls:
				return
			}
		}
	}()

	wg.Add(3)
	go func() {
		defer wg.Done()
		ch <- s.LastFM.Similar(0, src, opt.Lf...).Songs
	}()
	go func() {
		defer wg.Done()
		ch <- s.YaMusic.Similar(src, opt.Ya...).Songs
	}()
	go func() {
		defer wg.Done()
		ch <- s.Soundcloud.Similar(src, opt.Sc...).Songs
	}()

	wg.Wait()
	close(ch)
	cls <- struct{}{}
	close(cls)

	if opt.OnePerArtist {
		sort.SliceStable(res, func(i, j int) bool {
			return res[i].Artist < res[j].Artist
		})

		for i := 0; i < len(res)-1; i++ {
			if res[i].Artist == res[i+1].Artist {
				res = append(res[:i], res[i+1:]...)
				i--
			}
		}
	}

	return datastruct.Songs{
		Songs: res,
		From:  Frm,
	}
}

func (s WebApi) Top(artist string) datastruct.Songs {
	return s.LastFM.Top(strings.Split(artist, ", "), 10)
}

func (s WebApi) Close() {
	s.Soundcloud.Close()
}
