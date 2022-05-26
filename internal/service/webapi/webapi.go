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
)

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

func (s WebApiService) GetSimiliar(sourceData datastruct.AudioItems, oneAudioPerArtist bool) (result datastruct.AudioItems) {
	result = s.ILastFM.GetSimiliarSongsFromLast100(0, sourceData)
	result.Items = append(result.Items, s.IYaMusic.GetSimliarSongsFromYa100(sourceData).Items...)

	if oneAudioPerArtist {
		sort.SliceStable(result.Items, func(i, j int) bool {
			return result.Items[i].Artist < result.Items[j].Artist
		})

		for i, j := 0, 1; j < len(result.Items); i, j = i+1, j+1 {
			if result.Items[i].Artist == result.Items[j].Artist {
				result.Items = append(result.Items[:i], result.Items[j:]...)
			}
		}
	}

	return result
}

func (s WebApiService) GetTopSongs(artist string) datastruct.AudioItems {
	return s.ILastFM.GetTopTracks(strings.Split(artist, ", "), 10)
}

func (s WebApiService) TGSend() func(chattable tgbotapi.Chattable) tgbotapi.Message {
	return s.ITelegram.Send
}
