package yaMusic

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
)

const From = "yaMusic"

type YaMusic struct {
	log customLogger.Logger
	parser
}

func New(log customLogger.Logger) YaMusic {
	return YaMusic{
		parser: newParser(log),
		log:    log,
	}
}

func (y YaMusic) Find(query string) (audio datastruct.Song) {
	return y.parser.find(query)
}

func (y YaMusic) Similar(src datastruct.Songs, opts ...Set) datastruct.Songs {
	if opts != nil {
		return newClt(y.parser, opts...).similarParallel(src)
	}
	return newClt(y.parser).similarParallel(src)
}
