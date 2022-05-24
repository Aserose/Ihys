package youTube

import (
	"IhysBestowal/pkg/customLogger"
	ytsearch "github.com/AnjanaMadu/YTSearch"
)

type IYouTube interface {
	GetYTUrl(query string) string
}

type youTube struct {
	log customLogger.Logger
}

func NewDrafter(log customLogger.Logger) IYouTube {
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
			return "youtube.com/watch?v=" + results[0].VideoId
		}
		if len(results) > 1 {
			if results[1].VideoId != "" {
				return "youtube.com/watch?v=" + results[1].VideoId
			}
		}
	}

	return " "
}
