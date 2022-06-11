package soundcloud

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
)

type ISoundcloud interface {
	GetSimilar(sourceAudios datastruct.AudioItems, opts ...ProcessingOptions) datastruct.AudioItems
	CloseBrowser()
}

type soundcloud struct {
	log customLogger.Logger
	parser
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
