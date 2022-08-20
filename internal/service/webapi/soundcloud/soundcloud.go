package soundcloud

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
)

const From = "soundCloud"

type Soundcloud struct {
	parser
	log customLogger.Logger
}

func New(log customLogger.Logger) Soundcloud {
	return Soundcloud{
		parser: newParser(log),
		log:    log,
	}
}

func (s Soundcloud) Similar(src datastruct.Set, opts ...Set) datastruct.Set {
	if opts != nil {
		return newClt(s.parser, opts...).similarParallel(src)
	}
	return newClt(s.parser).similarParallel(src)
}

func (s Soundcloud) Close() {
	s.parser.CloseBrowser()
}
