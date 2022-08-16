package discogs

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
)

const empty = ``

type Discogs struct {
	clt
	enq
}

func New(log customLogger.Logger, cfg config.Discogs) Discogs {
	return Discogs{
		clt: newClt(),
		enq: newEnq(log, cfg),
	}
}

func (d Discogs) SiteLabel(query string) string {
	return d.clt.first(d.sites(query, empty))
}

func (d Discogs) SiteArtist(query string) string {
	return d.clt.first(d.sites(query, typeArtist))
}

func (d Discogs) SongInfo(audio datastruct.Song) datastruct.SongInfo {
	return d.enq.songInfo(audio)
}
