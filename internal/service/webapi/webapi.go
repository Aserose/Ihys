package webapi

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/repository"
	"IhysBestowal/internal/service/auth"
	"IhysBestowal/internal/service/webapi/lastFm"
	"IhysBestowal/internal/service/webapi/soundcloud"
	tgs "IhysBestowal/internal/service/webapi/tg"
	"IhysBestowal/internal/service/webapi/vk"
	"IhysBestowal/internal/service/webapi/yaMusic"
	"IhysBestowal/internal/service/webapi/youTube"
	"IhysBestowal/pkg/customLogger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sort"
	"strings"
	"sync"
)

type WebApiService struct {
	vk.IVk
	tgs.ITelegram
	youTube.IYouTube
	lastFm.ILastFM
	yaMusic.IYaMusic
	soundcloud.ISoundcloud
}

func NewWebApiService(log customLogger.Logger, cfg config.Service, repo repository.Repository, authService auth.AuthService) WebApiService {
	return WebApiService{
		ITelegram:   tgs.NewTg(log, cfg),
		IVk:         vk.NewVK(log, cfg.Vk, authService.Vk()),
		IYouTube:    youTube.NewYoutube(log),
		ILastFM:     lastFm.NewLastFM(log, cfg.LastFM, repo),
		IYaMusic:    yaMusic.NewYaMusic(log),
		ISoundcloud: soundcloud.NewSoundcloud(log),
	}
}

func (s WebApiService) Search(query string) datastruct.AudioItem {
	return s.IYaMusic.GetAudio(query)
}

func (s WebApiService) GetSimilar(sourceData datastruct.AudioItems, opt Opt) datastruct.AudioItems {
	wg := &sync.WaitGroup{}
	items := []datastruct.AudioItem{}
	ch := make(chan []datastruct.AudioItem)
	closed := make(chan bool)

	go func() {
		for {
			select {
			case i := <-ch:
				items = append(items, i...)
			case <-closed:
				return
			}
		}
	}()

	wg.Add(3)
	go func() {
		defer wg.Done()
		ch <- s.ILastFM.GetSimilar(0, sourceData, opt.lf...).Items
	}()
	go func() {
		defer wg.Done()
		ch <- s.IYaMusic.GetSimilar(sourceData, opt.ya...).Items
	}()
	go func() {
		defer wg.Done()
		ch <- s.ISoundcloud.GetSimilar(sourceData, opt.sc...).Items
	}()
	wg.Wait()

	close(ch)
	closed <- true
	close(closed)

	if opt.oneAudioPerArtist {
		sort.SliceStable(items, func(i, j int) bool {
			return items[i].Artist < items[j].Artist
		})

		for i := 0; i < len(items)-1; i++ {
			if items[i].Artist == items[i+1].Artist {
				items = append(items[:i], items[i+1:]...)
				i--
			}
		}
	}

	return datastruct.AudioItems{
		Items: items,
		From:  "all",
	}
}

func (s WebApiService) GetTopSongs(artist string) datastruct.AudioItems {
	return s.ILastFM.GetTopTracks(strings.Split(artist, ", "), 10)
}

func (s WebApiService) TGSend() func(chattable tgbotapi.Chattable) tgbotapi.Message {
	return s.ITelegram.Send
}
