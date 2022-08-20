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

func (md middleware) find(query string) datastruct.Song {
	return md.api.Find(query)
}

func (md middleware) similar(src datastruct.Songs) string {
	switch src.From {

	case md.api.From.All():
		return md.All(src)

	case md.api.From.YaMusic():
		return md.YaMusic(src)

	case md.api.From.Lfm().Similar():
		return md.LastFM(src)

	case md.api.From.Lfm().Top():
		return md.LastFMTop(src)

	}

	return emp
}

func (md middleware) All(src datastruct.Songs) string {
	if c := md.cache(src); c != emp {
		return c
	}
	return md.storage.Put(src.Songs[0], md.api.Similar(src, webapi.Default()))
}

func (md middleware) YaMusic(src datastruct.Songs) string {
	if c := md.cache(src); c != emp {
		return c
	}
	return md.storage.Put(src.Songs[0], md.api.YaMusic.Similar(src))
}

func (md middleware) LastFM(src datastruct.Songs) string {
	if c := md.cache(src); c != emp {
		return c
	}
	return md.storage.Put(src.Songs[0], md.api.LastFM.Similar(0, src))
}

func (md middleware) LastFMTop(src datastruct.Songs) string {
	if c := md.cache(src); c != emp {
		return c
	}
	return md.storage.Put(src.Songs[0], md.api.LastFM.Top([]string{src.Songs[0].Artist}, 7))
}

func (md middleware) VK(p dto.Response) datastruct.Songs {
	src, err := md.api.VK.Recommendations(p.TGUser, 0)
	if err != nil {
		md.api.TG.Send(tgbotapi.NewEditMessageText(p.ChatId, p.MsgId, err.Error()))
		return src
	}
	return src
}

func (md middleware) VKPlaylists(p dto.Response) datastruct.Playlists {
	src, err := md.api.VK.UserPlaylists(p.TGUser)
	if err != nil {
		md.api.TG.Send(tgbotapi.NewEditMessageText(p.ChatId, p.MsgId, err.Error()))
		return src
	}
	return src
}

func (md middleware) from() webapi.From {
	return md.api.From
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
	res := src.WithFrom(0)
	if i.storage.IsExist(res) {
		return res
	}
	return emp
}

func (i items) put(src datastruct.Song, similar datastruct.Songs) string {
	return i.storage.Put(src, similar)
}

func (i items) get(src string, page int) []datastruct.Song {
	if !i.storage.IsExist(src) {
		s := datastruct.Song{}.NewSongs(src)
		i.storage.Put(s.Songs[0], i.api.Similar(s, webapi.Default()))
	}
	return i.storage.Get(src, page)
}

func (i items) pageCount(src string) int {
	if !i.storage.IsExist(src) {
		s := datastruct.Song{}.NewSongs(src)
		i.storage.Put(s.Songs[0], i.api.Similar(s, webapi.Default()))
	}
	return i.storage.PageCount(src)
}

func (i items) pageCapacity() int {
	return i.storage.PageCapacity()
}
