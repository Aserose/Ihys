package youTube

import (
	"IhysBestowal/pkg/customLogger"
	ytsearch "github.com/AnjanaMadu/YTSearch"
)

const (
	videoUrl = "youtube.com/watch?v="
)

type IYouTube interface {
	GetYTUrl(query string) string
}

type youTube struct {
	log customLogger.Logger
}

func NewYoutube(log customLogger.Logger) IYouTube {
	return youTube{
		log: log,
	}
}

func (yt youTube) GetYTUrl(query string) string {
	results, err := ytsearch.Search(query)
	if err != nil {
		yt.log.Info(yt.log.CallInfoStr(), err.Error())
		return " "
	}

	if len(results) != 0 {
		if results[0].VideoId != "" {
			return videoUrl + results[0].VideoId
		}
		if len(results) > 1 {
			if results[1].VideoId != "" {
				return videoUrl + results[1].VideoId
			}
		}
	}

	return " "
}
