package lastFm

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/repository"
	"IhysBestowal/pkg/customLogger"
	"github.com/goccy/go-json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

const (
	qMethod      = `method`
	qLimit       = `limit`
	qTrack       = `track`
	qKey         = `api_key`
	qFormat      = `format`
	qArtist      = `artist`
	qAutocorrect = `autocorrect`
)

type enq struct {
	apiKey     string
	httpClient *http.Client
	log        customLogger.Logger
}

func newEnq(log customLogger.Logger, cfg config.LastFM, repo repository.Repository) enq {
	return enq{
		apiKey:     cfg.Key,
		httpClient: &http.Client{},
		log:        log,
	}
}

func (l enq) do(req *http.Request) []byte {
	resp, err := l.httpClient.Do(req)
	if err != nil {
		l.log.Warn(l.log.CallInfoStr(), err.Error())
	}
	b, _ := io.ReadAll(resp.Body)

	return b
}

func (l enq) find(query string) datastruct.Song {
	if query == emp {
		return datastruct.Song{}
	}

	resp := datastruct.LFMSearchTrack{}
	req, _ := http.NewRequest(http.MethodGet, bUrl, nil)
	req.URL.RawQuery = url.Values{
		qMethod: {mSearchTrack},
		qLimit:  {`1`},
		qTrack:  {query},
		qKey:    {l.apiKey},
		qFormat: {fJSON},
	}.Encode()

	json.Unmarshal(l.do(req), &resp)

	if len(resp.Results.TrackMatches.Tracks) == 0 {
		return datastruct.Song{}
	}

	return datastruct.Song{
		Artist: resp.Results.TrackMatches.Tracks[0].Artist,
		Title:  resp.Results.TrackMatches.Tracks[0].Name,
	}
}

func (l enq) similar(artist, title string) datastruct.Set {
	resp := datastruct.LFMUnmr{}
	req, _ := http.NewRequest(http.MethodGet, bUrl, nil)
	req.URL.RawQuery = url.Values{
		qMethod: {mGetSimilarTrack},
		qArtist: {artist},
		qTrack:  {url.QueryEscape(title)},
		qKey:    {l.apiKey},
		qFormat: {fJSON},
	}.Encode()

	json.Unmarshal(l.do(req), &resp)
	songs := make([]datastruct.Song, len(resp.LFMSimilarTracks.Tracks))

	for i, s := range resp.LFMSimilarTracks.Tracks {
		songs[i] = datastruct.Song{
			Artist: s.Artist.Name,
			Title:  s.Name,
		}
	}

	return datastruct.Set{
		Song: songs,
	}
}

func (l enq) top(artists []string, numPerArtist int) datastruct.Set {
	if artists == nil || numPerArtist < 1 {
		return datastruct.Set{}
	}

	res := make([]datastruct.Song, len(artists)*numPerArtist)
	wg := &sync.WaitGroup{}
	ch := make(chan datastruct.Song)
	cls := make(chan struct{})

	request := func(artist string) datastruct.LFMTopTracks {
		resp := datastruct.LFMUnmr{}
		req, _ := http.NewRequest(http.MethodGet, bUrl, nil)
		req.URL.RawQuery = url.Values{
			qMethod: {mGetTopTrack},
			qArtist: {artist},
			qKey:    {l.apiKey},
			qFormat: {fJSON},
		}.Encode()

		json.Unmarshal(l.do(req), &resp)

		return resp.LFMTopTracks
	}

	go func() {
		j := 0
		for {
			select {
			case track := <-ch:
				res[j] = track
				j++
			case <-cls:
				return
			}
		}
	}()

	for _, artist := range artists {
		wg.Add(1)
		go func(artist string) {
			defer wg.Done()

			for i, track := range request(artist).Tracks {
				if i >= numPerArtist {
					return
				}

				t := track

				ch <- datastruct.Song{Artist: t.Artist.Name, Title: t.Name}
			}
		}(artist)
	}
	wg.Wait()

	cls <- struct{}{}
	close(cls)
	close(ch)

	return datastruct.Set{
		From: FromTop,
		Song: res,
	}
}

func (l enq) similarArtists(artist string, max int) []string {
	if max <= 0 {
		return []string{}
	}

	res := []string{}
	wg := &sync.WaitGroup{}
	ch := make(chan []string)
	cls := make(chan struct{})

	request := func(artistName string) []string {
		resp := datastruct.LFMUnmr{}
		req, _ := http.NewRequest(http.MethodGet, bUrl, nil)
		req.URL.RawQuery = url.Values{
			qMethod:      {mGetSimilarArtist},
			qLimit:       {strconv.Itoa(max)},
			qArtist:      {artistName},
			qKey:         {l.apiKey},
			qFormat:      {fJSON},
			qAutocorrect: {`1`},
		}.Encode()

		json.Unmarshal(l.do(req), &resp)

		if resp.LFMSimilarArtists.Artists == nil {
			return []string{}
		}
		artists := make([]string, len(resp.LFMSimilarArtists.Artists))
		for i, r := range resp.LFMSimilarArtists.Artists {
			artists[i] = r.Name
		}
		return artists
	}

	go func() {
		for {
			select {
			case artists := <-ch:
				res = append(res, artists...)
			case <-cls:
				return
			}
		}
	}()

	if !func() (isEnum bool) {
		for _, enumType := range enumTypes {
			if strings.Contains(artist, enumType) {

				for _, name := range strings.Split(artist, enumType) {
					wg.Add(1)
					go func(name string) {
						defer wg.Done()
						ch <- request(name)
					}(name)
				}
				wg.Wait()
				if !isEnum {
					isEnum = true
				}

			}
		}
		return
	}() {
		ch <- request(artist)
	}

	cls <- struct{}{}
	close(cls)
	close(ch)

	if res == nil {
		return []string{}
	}

	return res
}

var enumTypes = []string{
	`, `,
	` feat. `,
	` feat `}
