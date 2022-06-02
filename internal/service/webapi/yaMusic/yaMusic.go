package yaMusic

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
)

type IYaMusic interface {
	GetSimilarSongsFromYa(sourceAudios datastruct.AudioItems, opts ...ProcessingOptions) datastruct.AudioItems
	GetAudio(query string) (audio datastruct.AudioItem)
}

type yaMusic struct {
	log customLogger.Logger
	parser
}

func NewYaMusic(log customLogger.Logger) IYaMusic {
	return yaMusic{
		parser: newParser(log),
		log:    log,
	}
}

func (y yaMusic) GetAudio(query string) (audio datastruct.AudioItem) {
	return y.parser.getAudio(query)
}

func (y yaMusic) GetSimilarSongsFromYa(sourceAudios datastruct.AudioItems, opts ...ProcessingOptions) datastruct.AudioItems {
	if opts != nil {
		return newCollater(y.parser, opts...).getSimilarParallel(sourceAudios)
	}
	return newCollater(y.parser).getSimilarParallel(sourceAudios)
}
