package menu

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/repository"
	"IhysBestowal/internal/service/webapi"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type middleware struct {
	menu webapi.Menu
	items
}

func newMiddleware(api webapi.WebApi, cache repository.Cache) middleware {
	return middleware{
		menu:  api.TG,
		items: newItems(cache, api),
	}
}

func (ms middleware) find(query string) datastruct.Song {
	return ms.api.Find(query)
}

func (ms middleware) similar(src datastruct.Songs) string {
	switch src.From {

	case ms.api.From.All():
		return ms.All(src)

	case ms.api.From.YaMusic():
		return ms.YaMusic(src)

	case ms.api.From.Lfm().Similar():
		return ms.LastFM(src)

	case ms.api.From.Lfm().Top():
		return ms.LastFMTop(src)
	}

	return emp
}

func (ms middleware) All(src datastruct.Songs) string {
	if cache := ms.cache(src); cache != emp {
		return cache
	}
	return ms.storage.Put(src.Songs[0], ms.api.Similar(src, webapi.Default()))
}

func (ms middleware) YaMusic(src datastruct.Songs) string {
	if cache := ms.cache(src); cache != emp {
		return cache
	}
	return ms.storage.Put(src.Songs[0], ms.api.YaMusic.Similar(src))
}

func (ms middleware) LastFM(src datastruct.Songs) string {
	if cache := ms.cache(src); cache != emp {
		return cache
	}
	return ms.storage.Put(src.Songs[0], ms.api.LastFM.Similar(0, src))
}

func (ms middleware) LastFMTop(src datastruct.Songs) string {
	if cache := ms.cache(src); cache != emp {
		return cache
	}
	return ms.storage.Put(src.Songs[0], ms.api.LastFM.Top([]string{src.Songs[0].Artist}, 7))
}

func (ms middleware) VK(p dto.Response) datastruct.Songs {
	src, err := ms.api.VK.Recommendations(p.TGUser, 0)
	if err != nil {
		ms.api.TG.Send(tgbotapi.NewEditMessageText(p.ChatId, p.MsgId, err.Error()))
		return src
	}

	return src
}

func (ms middleware) VKPlaylists(p dto.Response) datastruct.Playlists {
	src, err := ms.api.VK.UserPlaylists(p.TGUser)
	if err != nil {
		ms.api.TG.Send(tgbotapi.NewEditMessageText(p.ChatId, p.MsgId, err.Error()))
		return src
	}

	return src
}

func (ms middleware) from() webapi.From {
	return ms.api.From
}

type items struct {
	api     webapi.WebApi
	storage repository.Cache
}

func newItems(cache repository.Cache, api webapi.WebApi) items {
	return items{
		api:     api,
		storage: cache,
	}
}

func (i items) cache(src datastruct.Songs) string {
	sourceAudio := src.WithFrom(0)
	if i.storage.IsExist(sourceAudio) {
		return sourceAudio
	}
	return emp
}

func (i items) put(src datastruct.Song, similar datastruct.Songs) string {
	return i.storage.Put(src, similar)
}

func (i items) get(src string, page int) []datastruct.Song {
	if !i.storage.IsExist(src) {
		source := datastruct.Song{}.NewSongs(src)
		i.storage.Put(source.Songs[0], i.api.Similar(source, webapi.Default()))
	}
	return i.storage.Get(src, page)
}

func (i items) pageCount(src string) int {
	if !i.storage.IsExist(src) {
		source := datastruct.Song{}.NewSongs(src)
		i.storage.Put(source.Songs[0], i.api.Similar(source, webapi.Default()))
	}
	return i.storage.PageCount(src)
}

func (i items) pageCapacity() int {
	return i.storage.PageCapacity()
}
