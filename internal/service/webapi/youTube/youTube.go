package youTube

import (
	"IhysBestowal/pkg/customLogger"
	ytsearch "github.com/AnjanaMadu/YTSearch"
)

const (
	empty    = ``
	videoUrl = "youtube.com/watch?v="
)

type YouTube struct {
	log customLogger.Logger
}

func New(log customLogger.Logger) YouTube {
	return YouTube{
		log: log,
	}
}

func (yt YouTube) VideoURL(query string) string {
	res, err := ytsearch.Search(query)
	if err != nil {
		yt.log.Info(yt.log.CallInfoStr(), err.Error())
		return empty
	}

	if len(res) != 0 {
		if res[0].VideoId != empty {
			return videoUrl + res[0].VideoId
		}
		if len(res) > 1 {
			if res[1].VideoId != empty {
				return videoUrl + res[1].VideoId
			}
		}
	}

	return empty
}
