package yaMusic

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ndrewnee/go-yamusic/yamusic"
	"github.com/v-grabko1999/go-html2json"
	"io"
	"net/http"
	"strings"
)

type iParser interface {
	getSimiliars(artist, songTitle string) datastruct.YaMSimiliar
	getAudio(query string) (audio datastruct.AudioItem)
}

type parser struct {
	search func(query string) *yamusic.SearchResp
	log    customLogger.Logger
}

func newParser(log customLogger.Logger) iParser {
	p := parser{
		search: func(query string) *yamusic.SearchResp {
			respYa, _, _ := yamusic.NewClient(yamusic.HTTPClient(&http.Client{})).Search().Tracks(context.Background(), query, nil)
			return respYa
		},
		log: log,
	}

	return p
}

func (e parser) getSimiliars(artist, songTitle string) datastruct.YaMSimiliar {
	sourceData := e.getSidebarData(artist + " " + songTitle)

	if sourceData == nil {
		return datastruct.YaMSimiliar{}
	}

	return e.decode(sourceData)
}

func (e parser) decode(sourceData []byte) datastruct.YaMSimiliar {
	var (
		r   []datastruct.YaMSourcePage
		yaS datastruct.YaMSimiliar
	)

	json.Unmarshal(e.reformat(string(sourceData)), &r)
	json.Unmarshal([]byte(strings.TrimRight(strings.Trim(r[0].Elements[0].Elements[1].Elements[0].Text, "var Mu="), ";")), &yaS)

	return yaS
}

func (e parser) reformat(body string) []byte {
	d, err := html2json.New(strings.NewReader(body))
	if err != nil {
		e.log.Error(e.log.CallInfoStr(), err.Error())
	}

	j, err := d.ToJSON()
	if err != nil {
		e.log.Error(e.log.CallInfoStr(), err.Error())
	}

	return j
}

func (e parser) getAudio(query string) (audio datastruct.AudioItem) {
	result := e.search(query)

	if len(result.Result.Tracks.Results) != 0 {
		audio.Title = result.Result.Tracks.Results[0].Title
		for i, artist := range result.Result.Tracks.Results[0].Artists {
			audio.Artist += artist.Name
			if i != len(result.Result.Tracks.Results[0].Artists)-1 {
				audio.Artist += ", "
			}
		}
	}

	return audio
}

func (e parser) getSidebarData(query string) []byte {
	searchResult := e.search(query)

	if len(searchResult.Result.Tracks.Results) == 0 {
		return nil
	}

	resp, err := http.Get(fmt.Sprintf("https://music.yandex.ru/album/%d/track/%d",
		searchResult.Result.Tracks.Results[0].Albums[0].ID,
		searchResult.Result.Tracks.Results[0].ID))
	if err != nil {
		e.log.Error(e.log.CallInfoStr(), err.Error())
	}

	sourceData, _ := io.ReadAll(resp.Body)

	return sourceData
}
