package lastFm

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/repository"
	"IhysBestowal/pkg/customLogger"
)

const (
	From    = "lastFm"
	FromTop = "top"

	bUrl = "https://ws.audioscrobbler.com/2.0/"

	mGetSimilarArtist = "artist.getsimilar"
	mGetTopTrack      = "artist.gettoptracks"
	mGetSimilarTrack  = "track.getsimilar"
	mSearchTrack      = "track.search"

	fJSON = "json"

	emp = ``
)

type LastFM struct {
	enq
}

func New(log customLogger.Logger, cfg config.LastFM, repo repository.Repository) LastFM {
	return LastFM{
		enq: newEnq(log, cfg, repo),
	}
}

func (l LastFM) Auth(uid int64) {

}

func (l LastFM) Find(query string) datastruct.Song {
	return l.enq.find(query)
}

func (l LastFM) Top(artists []string, max int) datastruct.Songs {
	return l.enq.top(artists, max)
}

func (l LastFM) Similar(uid int64, src datastruct.Songs, opts ...Set) datastruct.Songs {
	return newClt(l.enq, opts...).SimilarParallel(uid, src)
}
