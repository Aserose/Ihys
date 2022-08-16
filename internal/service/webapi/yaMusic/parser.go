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

func (e parser) similar(artist, song string) []datastruct.Song {
	data := e.sidebarData(artist + " " + song)
	if data == nil {
		return []datastruct.Song{}
	}
	tracks := e.decode(data).SimilarTracks
	if tracks == nil {
		return []datastruct.Song{}
	}
	result := make([]datastruct.Song, len(tracks))

	for i, track := range tracks {
		result[i] = datastruct.Song{
			Title:  track.Title,
			Artist: e.name(track.Artists),
		}
	}

	return result
}

func (e parser) decode(d []byte) datastruct.YaMSimilar {
	data := []datastruct.YaMSourcePage{}
	res := datastruct.YaMSimilar{}

	json.Unmarshal(e.reformat(string(d)), &data)
	if data[0].Elements[0].Elements[1].Elements == nil {
		return res
	}
	json.Unmarshal([]byte(strings.TrimRight(strings.Trim(data[0].Elements[0].Elements[1].Elements[0].Text, "var Mu="), ";")), &res)

	return res
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

func (e parser) find(query string) (song datastruct.Song) {
	res := e.search(query)

	if len(res.Result.Tracks.Results) != 0 {
		song.Title = res.Result.Tracks.Results[0].Title
		for i, artist := range res.Result.Tracks.Results[0].Artists {
			song.Artist += artist.Name
			if i != len(res.Result.Tracks.Results[0].Artists)-1 {
				song.Artist += ", "
			}
		}
	} else {
		song.Artist = artistDefault
		song.Title = songTitleDefault
	}

	return
}

func (e parser) sidebarData(query string) []byte {
	s := e.search(query)

	if len(s.Result.Tracks.Results) == 0 {
		return nil
	}

	resp, err := http.Get(fmt.Sprintf(trackLink,
		s.Result.Tracks.Results[0].Albums[0].ID,
		s.Result.Tracks.Results[0].ID))
	if err != nil {
		e.log.Error(e.log.CallInfoStr(), err.Error())
	}

	data, _ := io.ReadAll(resp.Body)

	return data
}

func (e parser) name(artists []datastruct.YaMArtists) (res string) {
	if len(artists) > 1 {
		for i, artist := range artists {
			res += artist.Name
			if i < len(artists)-1 {
				res += ", "
			}
		}
	} else {
		res += artists[0].Name
	}

	return
}
