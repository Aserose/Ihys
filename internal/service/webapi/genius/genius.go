package genius

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
)

type Genius struct {
	enq enq
}

func New(log customLogger.Logger, cfg config.Genius) Genius {
	return Genius{
		enq: newEnq(log, cfg),
	}
}

func (g Genius) LyricsURL(audio datastruct.Song) string {
	return g.enq.lyricsURL(audio)
}
