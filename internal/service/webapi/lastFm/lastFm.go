package lastFm

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/repository"
	"IhysBestowal/pkg/customLogger"
)

const (
	SourceFrom    = "lastFm"
	SourceFromTop = "top"

	baseUrl = "https://ws.audioscrobbler.com/2.0/?"

	methodGetSimilarArtist = "artist.getsimilar"
	methodGetTopTrack      = "artist.gettoptracks"
	methodGetSimilarTrack  = "track.getsimilar"
	methodSearchTrack      = "track.search"

	formatJSON             = "json"

	empty = ``
)

type ILastFM interface {
	Authorize(userId int64)
	GetSimilar(userId int64, sourceData datastruct.AudioItems, opts ...ProcessingOptions) datastruct.AudioItems
	GetTopTracks(artistNames []string, numberOfSongs int) datastruct.AudioItems
	GetAudio(query string) datastruct.AudioItem
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

func (l lastFm) GetAudio(query string) datastruct.AudioItem {
	return l.enquirer.getAudio(query)
}

func (l lastFm) GetSimilar(userId int64, sourceData datastruct.AudioItems, opts ...ProcessingOptions) datastruct.AudioItems {
	return newCollater(l.enquirer, opts...).getSimilarParallel(userId, sourceData)
}

func (l lastFm) GetTopTracks(artistNames []string, numberOfSongs int) datastruct.AudioItems {
	return l.enquirer.getTopTracks(artistNames, numberOfSongs)
}
