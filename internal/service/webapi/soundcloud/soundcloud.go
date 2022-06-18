package soundcloud

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
)

const SourceFrom = "soundCloud"

type ISoundcloud interface {
	GetSimilar(sourceAudios datastruct.AudioItems, opts ...ProcessingOptions) datastruct.AudioItems
	Close()
}

type soundcloud struct {
	parser
	log customLogger.Logger
}

func NewSoundcloud(log customLogger.Logger) ISoundcloud {
	return soundcloud{
		parser: newParser(log),
		log:    log,
	}
}

func (s soundcloud) GetSimilar(sourceAudios datastruct.AudioItems, opts ...ProcessingOptions) datastruct.AudioItems {
	if opts != nil {
		return newCollater(s.parser, opts...).getSimilarParallel(sourceAudios)
	}
	return newCollater(s.parser).getSimilarParallel(sourceAudios)
}

func (s soundcloud) Close() {
	s.parser.CloseBrowser()
}
