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
	queryMethod      = `method`
	queryLimit       = `limit`
	queryTrack       = `track`
	queryKey         = `api_key`
	queryFormat      = `format`
	queryArtist      = `artist`
	queryAutocorrect = `autocorrect`
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

func (l enq) request(req *http.Request) []byte {
	resp, err := l.httpClient.Do(req)
	if err != nil {
		l.log.Warn(l.log.CallInfoStr(), err.Error())
	}
	b, _ := io.ReadAll(resp.Body)

	return b
}

func (l enq) find(query string) datastruct.Song {
	if query == empty {
		return datastruct.Song{}
	}

	resp := datastruct.LastFMSearchTrackResult{}
	req, _ := http.NewRequest(http.MethodGet, baseUrl, nil)
	req.URL.RawQuery = url.Values{
		queryMethod: {methodSearchTrack},
		queryLimit:  {`1`},
		queryTrack:  {query},
		queryKey:    {l.apiKey},
		queryFormat: {formatJSON},
	}.Encode()

	json.Unmarshal(l.request(req), &resp)

	if len(resp.Results.TrackMatches.Tracks) == 0 {
		return datastruct.Song{}
	}

	return datastruct.Song{
		Artist: resp.Results.TrackMatches.Tracks[0].Artist,
		Title:  resp.Results.TrackMatches.Tracks[0].Name,
	}
}

func (l enq) similar(artist, title string) datastruct.Songs {
	resp := datastruct.LastFMUnmr{}
	req, _ := http.NewRequest(http.MethodGet, baseUrl, nil)
	req.URL.RawQuery = url.Values{
		queryMethod: {methodGetSimilarTrack},
		queryArtist: {artist},
		queryTrack:  {url.QueryEscape(title)},
		queryKey:    {l.apiKey},
		queryFormat: {formatJSON},
	}.Encode()

	json.Unmarshal(l.request(req), &resp)
	songs := make([]datastruct.Song, len(resp.LastFMSimilarTracks.Tracks))

	for i, s := range resp.LastFMSimilarTracks.Tracks {
		songs[i] = datastruct.Song{
			Artist: s.Artist.Name,
			Title:  s.Name,
		}
	}

	return datastruct.Songs{
		Songs: songs,
	}
}

func (l enq) top(artists []string, numPerArtist int) datastruct.Songs {
	if artists == nil || numPerArtist <= 0 {
		return datastruct.Songs{}
	}

	res := make([]datastruct.Song, len(artists)*numPerArtist)
	wg := &sync.WaitGroup{}
	ch := make(chan datastruct.Song)
	cls := make(chan struct{})

	request := func(artist string) datastruct.LastFMTopTracks {
		resp := datastruct.LastFMUnmr{}
		req, _ := http.NewRequest(http.MethodGet, baseUrl, nil)
		req.URL.RawQuery = url.Values{
			queryMethod: {methodGetTopTrack},
			queryArtist: {artist},
			queryKey:    {l.apiKey},
			queryFormat: {formatJSON},
		}.Encode()

		json.Unmarshal(l.request(req), &resp)

		return resp.LastFMTopTracks
	}

	go func() {
		j := 0
		for {
			select {
			case track, ok := <-ch:
				if !ok {
					continue
				}
				res[j] = track
				j++
			case <-cls:
				return
			}
		}
	}()

	wg.Add(len(artists))
	for _, artist := range artists {
		go func(artist string) {
			defer wg.Done()

			for i, track := range request(artist).Tracks {
				t := track
				if i >= numPerArtist {
					break
				}
				ch <- datastruct.Song{
					Artist: t.Artist.Name,
					Title:  t.Name,
				}
			}
		}(artist)
	}
	wg.Wait()
	close(ch)
	cls <- struct{}{}
	close(cls)

	return datastruct.Songs{
		Songs: res,
		From:  FromTop,
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
		resp := datastruct.LastFMUnmr{}
		req, _ := http.NewRequest(http.MethodGet, baseUrl, nil)
		req.URL.RawQuery = url.Values{
			queryMethod:      {methodGetSimilarArtist},
			queryLimit:       {strconv.Itoa(max)},
			queryArtist:      {artistName},
			queryKey:         {l.apiKey},
			queryFormat:      {formatJSON},
			queryAutocorrect: {`1`},
		}.Encode()

		json.Unmarshal(l.request(req), &resp)

		if resp.LastFMSimilarArtists.Artists == nil {
			return []string{}
		}
		artistList := make([]string, len(resp.LastFMSimilarArtists.Artists))
		for i, r := range resp.LastFMSimilarArtists.Artists {
			artistList[i] = r.Name
		}
		return artistList
	}

	go func() {
		for {
			select {
			case names, ok := <-ch:
				if !ok {
					continue
				}
				res = append(res, names...)
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
	close(ch)
	cls <- struct{}{}
	close(cls)

	if res == nil {
		return []string{}
	}

	return res
}

var enumTypes = []string{
	`, `,
	` feat. `,
	` feat `}
