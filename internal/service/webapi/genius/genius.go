package genius

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
)

type IGenius interface {
	GetLyricsURL(audio datastruct.AudioItem) string
}

type genius struct {
	enq enquirer
}

func NewGenius(log customLogger.Logger, cfg config.Genius) IGenius {
	return genius{
		enq: newEnquirer(log, cfg),
	}
}

func (g genius) GetLyricsURL(audio datastruct.AudioItem) string {
	return g.enq.getLyricsURL(audio)
}
