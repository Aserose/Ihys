package lastFm

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/repository"
	"IhysBestowal/pkg/customLogger"
)

const (
	baseUri          = "https://ws.audioscrobbler.com/2.0/?"
	getSimilarArtist = "artist.getsimilar"
	getSimilarTrack  = "track.getsimilar"
	getTopTrack      = "artist.gettoptracks"
	jsonFrmt         = "json"
)


type ILastFM interface {
	Authorize(userId int64)
	GetSimiliarSongsFromLast(userId int64, sourceData datastruct.AudioItems, opts ...ProcessingOptions) datastruct.AudioItems
	GetTopTracks(artistNames []string, numberOfSongs int) datastruct.AudioItems
}

type lastFm struct {
	enquirer
}

func NewLastFM(log customLogger.Logger, cfg config.LastFM, repo repository.Repository) ILastFM {
	return &lastFm{
		enquirer: newEnquirer(log, cfg, repo),
	}
}

func (l lastFm) Authorize(userId int64) {

	//TODO

}

func (l lastFm) GetSimiliarSongsFromLast(userId int64, sourceData datastruct.AudioItems, opts ...ProcessingOptions) datastruct.AudioItems {
	return newCollater(l.enquirer, opts...).getSimilarParallel(userId, sourceData)
}

func (l lastFm) GetTopTracks(artistNames []string, numberOfSongs int) datastruct.AudioItems {
	return l.enquirer.getTopTracks(artistNames, numberOfSongs)
}
