package menu

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/service/webapi"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type iMiddleware interface {
	GetVKPlaylists(tr dto.Executor, msgId int) datastruct.PlaylistItems
	GetYaMusicSimiliars(tr dto.Executor, msgId int) datastruct.AudioItems
	GetLastFMSimiliars(tr dto.Executor, msgId int) datastruct.AudioItems
	GetVKSimiliars(tr dto.Executor, msgId int) datastruct.AudioItems
}

type middleware struct {
	api    webapi.WebApiService
	source func(tr dto.Executor, msgId int) datastruct.AudioItems
}

func newMiddleware(api webapi.WebApiService, audioSource func(tr dto.Executor, msgId int) datastruct.AudioItems) iMiddleware {
	g := middleware{
		api: api,
	}

	if audioSource == nil {
		g.source = g.GetVKSimiliars
	} else {
		g.source = audioSource
	}

	return g
}

func (ms middleware) GetYaMusicSimiliars(tr dto.Executor, msgId int) datastruct.AudioItems {
	sourceData := ms.api.IYaMusic.GetSimliarSongsFromYa(ms.source(tr, msgId))

	sourceData.From = "YaMusic"

	return sourceData
}

func (ms middleware) GetLastFMSimiliars(tr dto.Executor, msgId int) datastruct.AudioItems {
	sourceData := ms.api.ILastFM.GetSimiliarSongsFromLast(0, ms.source(tr, msgId))

	sourceData.From = "LastFM"

	return sourceData
}

func (ms middleware) GetVKPlaylists(tr dto.Executor, msgId int) datastruct.PlaylistItems {
	sourceData, err := ms.api.IVk.GetUserPlaylists(tr.TGUser)
	if err != nil {
		ms.api.Send(tgbotapi.NewEditMessageText(tr.ChatID, msgId, err.Error()))
		return sourceData
	}

	sourceData.From = "VK"

	return sourceData
}

func (ms middleware) GetVKSimiliars(tr dto.Executor, msgId int) datastruct.AudioItems {
	sourceData, err := ms.api.IVk.GetRecommendations(tr.TGUser, 0)
	if err != nil {
		ms.api.Send(tgbotapi.NewEditMessageText(tr.ChatID, msgId, err.Error()))
		return sourceData
	}

	sourceData.From = "VK"

	return sourceData
}
