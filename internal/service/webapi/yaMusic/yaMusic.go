package yaMusic

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
)

type IYaMusic interface {
	GetSimliarSongsFromYa(data datastruct.AudioItems) datastruct.AudioItems
	GetSimliarSongsFromYa100(sourceAudios datastruct.AudioItems) datastruct.AudioItems
	GetAudio(query string) (audio datastruct.AudioItem)
}

type yaMusic struct {
	log customLogger.Logger
}

func NewYaMusic(log customLogger.Logger) IYaMusic {
	return yaMusic{
		log: log,
	}
}

func (y yaMusic) GetAudio(query string) (audio datastruct.AudioItem) {
	return newEnquirer(y.log).getAudio(query)
}

func (y yaMusic) GetSimliarSongsFromYa100(sourceAudios datastruct.AudioItems) datastruct.AudioItems {
	return newEnquirer(y.log, setAudioAmountPerSource(100)).getSimiliar(sourceAudios)
}

func (y yaMusic) GetSimliarSongsFromYa(sourceAudios datastruct.AudioItems) datastruct.AudioItems {
	return newEnquirer(y.log).getSimiliar(sourceAudios)
}
