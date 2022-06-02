package webapi

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/repository"
	"IhysBestowal/internal/service/auth"
	"IhysBestowal/internal/service/webapi/lastFm"
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

type Opt struct {
	ya []yaMusic.ProcessingOptions
	lf []lastFm.ProcessingOptions
}

type WebApiService struct {
	vk.IVk
	tgs.ITelegram
	youTube.IYouTube
	lastFm.ILastFM
	yaMusic.IYaMusic
}

func NewWebApiService(log customLogger.Logger, cfg config.Service, repo repository.Repository, authService auth.AuthService) WebApiService {
	return WebApiService{
		ITelegram: tgs.NewTg(log, cfg),
		IVk:       vk.NewVK(log, cfg.Vk, authService.Vk()),
		IYouTube:  youTube.NewYoutube(log),
		ILastFM:   lastFm.NewLastFM(log, cfg.LastFM, repo),
		IYaMusic:  yaMusic.NewYaMusic(log),
	}
}

func (s WebApiService) Search(query string) datastruct.AudioItem {
	return s.IYaMusic.GetAudio(query)
}

func (s WebApiService) GetSimilar(sourceData datastruct.AudioItems, oneAudioPerArtist bool, opt Opt) datastruct.AudioItems {
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

	wg.Add(2)
	go func() {
		defer wg.Done()
		ch <- s.ILastFM.GetSimiliarSongsFromLast(0, sourceData, opt.lf...).Items
	}()
	go func() {
		defer wg.Done()
		ch <- s.IYaMusic.GetSimilarSongsFromYa(sourceData, opt.ya...).Items
	}()
	wg.Wait()

	close(ch)
	closed <- true

	if oneAudioPerArtist {
		sort.SliceStable(items, func(i, j int) bool {
			return items[i].Artist < items[j].Artist
		})

		for i, j := 0, 1; j < len(items); i, j = i+1, j+1 {
			if items[i].Artist == items[j].Artist {
				items = append(items[:i], items[j:]...)
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
