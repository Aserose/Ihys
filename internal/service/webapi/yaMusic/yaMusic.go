package yaMusic

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
)

type IYaMusic interface {
	GetSimliarSongsFromYa(data datastruct.AudioItems) datastruct.AudioItems
	GetSimliarSongsFromYa10(sourceAudios datastruct.AudioItems) datastruct.AudioItems
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
	return newCollater(y.log).getAudio(query)
}

// Getting the maximum number of similar songs, which value is 10 per audio source.
func (y yaMusic) GetSimliarSongsFromYa10(sourceAudios datastruct.AudioItems) datastruct.AudioItems {
	return newCollater(y.log, setMaxAudioAmountPerSource(10)).getSimiliar(sourceAudios)
}

// Getting 3 similar songs to the source.
func (y yaMusic) GetSimliarSongsFromYa(sourceAudios datastruct.AudioItems) datastruct.AudioItems {
	return newCollater(y.log).getSimiliar(sourceAudios)
}
