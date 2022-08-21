package menu

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/repository"
	"IhysBestowal/internal/service/webapi"
)

type TGMenu struct {
	md middleware
	viewer
	cfg config.Keypads
}

func New(api webapi.WebApi, cache repository.Cache, cfg config.Keypads) TGMenu {
	md := newMiddleware(api, cache)
	return TGMenu{
		md:     md,
		viewer: newViewer(cfg, md, api),
		cfg:    cfg,
	}
}

func (ms TGMenu) Preload(p dto.Response) {
	ms.viewer.preload(p)
}

func (ms TGMenu) Random(p dto.Response) {
	ms.openSong(p, datastruct.Set{
		Song: []datastruct.Song{ms.md.api.Random()},
	})
}

func (ms TGMenu) Find(p dto.Response, query string) {
	ms.openSong(p, datastruct.Set{
		Song: []datastruct.Song{ms.md.find(query)},
	})
}

func (ms TGMenu) openSong(p dto.Response, src datastruct.Set) {
	ms.viewer.openSong(p, src.Song[0])
}

func (ms TGMenu) Main(p dto.Response) {

	// TODO

}

// TODO AUTH MENU

// TODO PLAYLIST
