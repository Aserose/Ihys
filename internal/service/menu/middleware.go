package menu

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/repository"
	"IhysBestowal/internal/service/webapi"
	"IhysBestowal/internal/service/webapi/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type middleware struct {
	api       webapi.WebApiService
	tgBuilder tg.TGMenu
	items
}

func newMiddleware(api webapi.WebApiService, storage repository.TrackStorage) middleware {
	return middleware{
		api:       api,
		tgBuilder: api.ITelegram.NewMenuBuilder(),
		items:     newItems(storage),
	}
}

func (ms middleware) search(query string) datastruct.AudioItem {
	return ms.api.Search(query)
}

func (ms middleware) getAllSimilar(source datastruct.AudioItems) string {
	if sourceAudio := ms.isExist(source, ms.api.GetSourceFrom.All()); sourceAudio != empty {
		return sourceAudio
	}
	return ms.storage.Put(source.Items[0], ms.api.GetSimilar(source, webapi.GetOptDefaultPreset()))
}

func (ms middleware) getYaMusicSimilar(source datastruct.AudioItems) string {
	if sourceAudio := ms.isExist(source, ms.api.GetSourceFrom.YaMusic()); sourceAudio != empty {
		return sourceAudio
	}
	return ms.storage.Put(source.Items[0], ms.api.IYaMusic.GetSimilar(source))
}

func (ms middleware) getLastFMSimilar(source datastruct.AudioItems) string {
	if sourceAudio := ms.isExist(source, ms.api.GetSourceFrom.Lfm().LastFm); sourceAudio != empty {
		return sourceAudio
	}
	return ms.storage.Put(source.Items[0], ms.api.ILastFM.GetSimilar(0, source))
}

func (ms middleware) getLastFMBest(source datastruct.AudioItems) string {
	if sourceAudio := ms.isExist(source, ms.api.GetSourceFrom.Lfm().LastFmTop); sourceAudio != empty {
		return sourceAudio
	}
	return ms.storage.Put(source.Items[0], ms.api.ILastFM.GetTopTracks([]string{source.Items[0].Artist}, 7))
}

func (ms middleware) getVKSimilar(p dto.Response) datastruct.AudioItems {
	sourceData, err := ms.api.IVk.GetRecommendations(p.TGUser, 0)
	if err != nil {
		ms.api.Send(tgbotapi.NewEditMessageText(p.ChatId, p.MsgId, err.Error()))
		return sourceData
	}

	return sourceData
}

func (ms middleware) getVKPlaylists(p dto.Response) datastruct.PlaylistItems {
	sourceData, err := ms.api.IVk.GetUserPlaylists(p.TGUser)
	if err != nil {
		ms.api.Send(tgbotapi.NewEditMessageText(p.ChatId, p.MsgId, err.Error()))
		return sourceData
	}

	return sourceData
}

type items struct {
	storage repository.TrackStorage
}

func newItems(storage repository.TrackStorage) items {
	return items{
		storage: storage,
	}
}

func (i items) isExist(song datastruct.AudioItems, sourceFrom string) string {
	sourceAudio := song.Items[0].GetSourceAudio(sourceFrom)
	if i.storage.IsExist(sourceAudio) {
		return sourceAudio
	}
	return empty
}

func (i items) put(sourceAudio datastruct.AudioItem, similar datastruct.AudioItems) string {
	return i.storage.Put(sourceAudio, similar)
}

func (i items) get(sourceAudio string, page int) []datastruct.AudioItem {
	return i.storage.GetItems(sourceAudio, page)
}

func (i items) pageCount(sourceAudio string) int {
	return i.storage.GetPageCount(sourceAudio)
}
func (i items) pageCapacity() int {
	return i.storage.GetPageCapacity()
}
