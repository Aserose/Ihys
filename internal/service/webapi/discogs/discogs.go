package discogs

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
)

const (
	empty = ``
)

type IDiscogs interface {
	GetWebsiteLabel(query string) string
	GetWebsiteArtist(query string) string
	GetSongInfo(audio datastruct.AudioItem) datastruct.AudioInfo
}

type discogs struct {
	collater
	enquirer
}

func NewDiscogs(log customLogger.Logger, cfg config.Discogs) IDiscogs {
	return discogs{
		collater: newCollater(),
		enquirer: newEnquirer(log, cfg),
	}
}

func (d discogs) GetWebsiteLabel(query string) string {
	return d.collater.getFirstWebsite(d.getWebsites(query, empty))
}

func (d discogs) GetWebsiteArtist(query string) string {
	return d.collater.getFirstWebsite(d.getWebsites(query, typeArtist))
}

func (d discogs) GetSongInfo(audio datastruct.AudioItem) datastruct.AudioInfo {
	return d.enquirer.getSongInfo(audio)
}
