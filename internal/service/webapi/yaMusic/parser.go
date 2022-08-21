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

const trackLink = "https://music.yandex.ru/album/%d/track/%d"

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
	sidebar := e.sidebarData(artist + " " + song)
	if sidebar == nil {
		return []datastruct.Song{}
	}
	tracks := e.decode(sidebar).SimilarTracks
	if tracks == nil {
		return []datastruct.Song{}
	}

	res := make([]datastruct.Song, len(tracks))

	for i, track := range tracks {
		res[i] = datastruct.Song{
			Title:  track.Title,
			Artist: e.merge(track.Artists),
		}
	}

	return res
}

func (e parser) decode(d []byte) datastruct.YaSimilar {
	srcPage := []datastruct.YaSongPage{}
	res := datastruct.YaSimilar{}

	json.Unmarshal(e.json(string(d)), &srcPage)
	if srcPage[0].Elements[0].Elements[1].Elements == nil {
		return res
	}
	json.Unmarshal([]byte(strings.TrimRight(strings.Trim(srcPage[0].Elements[0].Elements[1].Elements[0].Text, "var Mu="), ";")), &res)

	return res
}

func (e parser) json(html string) []byte {
	d, err := html2json.New(strings.NewReader(html))
	if err != nil {
		e.log.Error(e.log.CallInfoStr(), err.Error())
	}

	j, err := d.ToJSON()
	if err != nil {
		e.log.Error(e.log.CallInfoStr(), err.Error())
	}

	return j
}

func (e parser) find(query string) (s datastruct.Song) {
	f := e.search(query)

	if len(f.Result.Tracks.Results) > 0 {
		s.Title = f.Result.Tracks.Results[0].Title

		for i, artist := range f.Result.Tracks.Results[0].Artists {
			s.Artist += artist.Name

			if i != len(f.Result.Tracks.Results[0].Artists)-1 {
				s.Artist += ", "
			}
		}

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

func (e parser) merge(artists []datastruct.YaArtists) (res string) {
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
