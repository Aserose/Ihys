package lastFm

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/repository"
	"IhysBestowal/pkg/customLogger"
)

type ILastFM interface {
	Authorize(userId int64)
	GetSimiliarSongsFromLast100(userId int64, sourceData datastruct.AudioItems) datastruct.AudioItems
	GetSimiliarSongsFromLast(userId int64, sourceData datastruct.AudioItems) datastruct.AudioItems
	GetTopTracks(artistNames []string, numberOfSongs int) datastruct.AudioItems
}

type lastFm struct {
	iEnquirer
}

func NewLastFM(log customLogger.Logger, cfg config.LastFM, repo repository.Repository) ILastFM {
	return &lastFm{
		iEnquirer: newEnquirer(log, cfg, repo),
	}
}

func (l lastFm) Authorize(userId int64) {

	//TODO

}

func (l lastFm) GetSimiliarSongsFromLast100(userId int64, sourceData datastruct.AudioItems) datastruct.AudioItems {
	return newCollater(l.iEnquirer, setMaxAudioAmountPerSource(100)).getSimiliars(userId, sourceData)
}

func (l lastFm) GetSimiliarSongsFromLast(userId int64, sourceData datastruct.AudioItems) datastruct.AudioItems {
	return newCollater(l.iEnquirer).getSimiliars(userId, sourceData)
}

func (l lastFm) GetTopTracks(artistNames []string, numberOfSongs int) datastruct.AudioItems {
	return l.iEnquirer.getTopTracks(artistNames, numberOfSongs)
}
