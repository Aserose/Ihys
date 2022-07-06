package yaMusic

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/ndrewnee/go-yamusic/yamusic"
	"github.com/v-grabko1999/go-html2json"
	"io"
	"net/http"
	"strings"
)

const (
	artistDefault    = `Rick Astley`
	songTitleDefault = `Never Gonna Give You Up`

	trackLink = "https://music.yandex.ru/album/%d/track/%d"
)

type parser struct {
	client *yamusic.Client
	log    customLogger.Logger
}

func newParser(log customLogger.Logger) parser {
	return parser{
		client: yamusic.NewClient(),
		log:    log,
	}
}

func (e parser) search(query string) *yamusic.SearchResp {
	resp, _, _ := e.client.Search().Tracks(context.Background(), query, nil)

	return resp
}

func (e parser) getSimilar(artist, songTitle string) []datastruct.AudioItem {
	sourceData := e.getSidebarData(artist + " " + songTitle)
	if sourceData == nil {
		return []datastruct.AudioItem{}
	}
	yaTracks := e.decode(sourceData).SimilarTracks
	if yaTracks == nil {
		return []datastruct.AudioItem{}
	}
	result := make([]datastruct.AudioItem, len(yaTracks))

	for i, track := range yaTracks {
		result[i] = datastruct.AudioItem{
			Title:  track.Title,
			Artist: e.writeArtistName(track.Artists),
		}
	}

	return result
}

func (e parser) decode(sourceData []byte) datastruct.YaMSimilar {
	r := []datastruct.YaMSourcePage{}
	yaS := datastruct.YaMSimilar{}

	json.Unmarshal(e.reformat(string(sourceData)), &r)
	if r[0].Elements[0].Elements[1].Elements == nil {
		return yaS
	}
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
	} else {
		audio.Artist = artistDefault
		audio.Title = songTitleDefault
	}

	return
}

func (e parser) getSidebarData(query string) []byte {
	searchResult := e.search(query)

	if len(searchResult.Result.Tracks.Results) == 0 {
		return nil
	}

	resp, err := http.Get(fmt.Sprintf(trackLink,
		searchResult.Result.Tracks.Results[0].Albums[0].ID,
		searchResult.Result.Tracks.Results[0].ID))
	if err != nil {
		e.log.Error(e.log.CallInfoStr(), err.Error())
	}

	sourceData, _ := io.ReadAll(resp.Body)

	return sourceData
}

func (e parser) writeArtistName(artists []datastruct.YaMArtists) (result string) {
	if len(artists) > 1 {
		for i, artist := range artists {
			result += artist.Name
			if i < len(artists)-1 {
				result += ", "
			}
		}
	} else {
		result += artists[0].Name
	}

	return
}
