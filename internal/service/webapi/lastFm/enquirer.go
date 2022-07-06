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

type enquirer struct {
	apiKey     string
	httpClient *http.Client
	log        customLogger.Logger
}

func newEnquirer(log customLogger.Logger, cfg config.LastFM, repo repository.Repository) enquirer {
	return enquirer{
		apiKey:     cfg.Key,
		httpClient: &http.Client{},
		log:        log,
	}
}

func (l enquirer) sendRequest(req *http.Request) []byte {
	resp, err := l.httpClient.Do(req)
	if err != nil {
		l.log.Warn(l.log.CallInfoStr(), err.Error())
	}
	b, _ := io.ReadAll(resp.Body)

	return b
}

func (l enquirer) getAudio(query string) datastruct.AudioItem {
	if query == empty {
		return datastruct.AudioItem{}
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

	json.Unmarshal(l.sendRequest(req), &resp)

	if len(resp.Results.TrackMatches.Tracks) == 0 {
		return datastruct.AudioItem{}
	}

	return datastruct.AudioItem{
		Artist: resp.Results.TrackMatches.Tracks[0].Artist,
		Title:  resp.Results.TrackMatches.Tracks[0].Name,
	}
}

func (l enquirer) getSimilarTracks(artist, track string) datastruct.AudioItems {
	resp := datastruct.LastFMUnmr{}
	req, _ := http.NewRequest(http.MethodGet, baseUrl, nil)
	req.URL.RawQuery = url.Values{
		queryMethod: {methodGetSimilarTrack},
		queryArtist: {artist},
		queryTrack:  {url.QueryEscape(track)},
		queryKey:    {l.apiKey},
		queryFormat: {formatJSON},
	}.Encode()

	json.Unmarshal(l.sendRequest(req), &resp)
	trackList := make([]datastruct.AudioItem, len(resp.LastFMSimilarTracks.Tracks))

	for i, s := range resp.LastFMSimilarTracks.Tracks {
		trackList[i] = datastruct.AudioItem{
			Artist: s.Artist.Name,
			Title:  s.Name,
		}
	}

	return datastruct.AudioItems{
		Items: trackList,
	}
}

func (l enquirer) getTopTracks(artistNames []string, numberOfTracksPerArtist int) datastruct.AudioItems {
	if artistNames == nil || numberOfTracksPerArtist <= 0 {
		return datastruct.AudioItems{}
	}

	trackList := make([]datastruct.AudioItem, len(artistNames)*numberOfTracksPerArtist)
	wg := &sync.WaitGroup{}
	ch := make(chan datastruct.AudioItem)
	closed := make(chan bool)

	request := func(artistName string) datastruct.LastFMTopTracks {
		resp := datastruct.LastFMUnmr{}
		req, _ := http.NewRequest(http.MethodGet, baseUrl, nil)
		req.URL.RawQuery = url.Values{
			queryMethod: {methodGetTopTrack},
			queryArtist: {artistName},
			queryKey:    {l.apiKey},
			queryFormat: {formatJSON},
		}.Encode()

		json.Unmarshal(l.sendRequest(req), &resp)

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
				trackList[j] = track
				j++
			case <-closed:
				return
			}
		}
	}()

	wg.Add(len(artistNames))
	for _, artistName := range artistNames {
		go func(artistName string) {
			defer wg.Done()

			for i, track := range request(artistName).Tracks {
				t := track
				if i >= numberOfTracksPerArtist {
					break
				}
				ch <- datastruct.AudioItem{
					Artist: t.Artist.Name,
					Title:  t.Name,
				}
			}
		}(artistName)
	}
	wg.Wait()
	close(ch)
	closed <- true
	close(closed)

	return datastruct.AudioItems{
		Items: trackList,
		From:  SourceFromTop,
	}
}

func (l enquirer) getSimilarArtists(artistName string, limit int) []string {
	if limit <= 0 {
		return []string{}
	}

	result := []string{}
	wg := &sync.WaitGroup{}
	ch := make(chan []string)
	closed := make(chan bool)

	request := func(artistName string) []string {
		resp := datastruct.LastFMUnmr{}
		req, _ := http.NewRequest(http.MethodGet, baseUrl, nil)
		req.URL.RawQuery = url.Values{
			queryMethod:      {methodGetSimilarArtist},
			queryLimit:       {strconv.Itoa(limit)},
			queryArtist:      {artistName},
			queryKey:         {l.apiKey},
			queryFormat:      {formatJSON},
			queryAutocorrect: {`1`},
		}.Encode()

		json.Unmarshal(l.sendRequest(req), &resp)

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
				result = append(result, names...)
			case <-closed:
				return
			}
		}
	}()

	if !func() (isEnum bool) {
		for _, enumType := range enumTypes {
			if strings.Contains(artistName, enumType) {

				for _, name := range strings.Split(artistName, enumType) {
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
		ch <- request(artistName)
	}
	close(ch)
	closed <- true
	close(closed)

	if result == nil {
		return []string{}
	}

	return result
}

var enumTypes = []string{
	`, `,
	` feat. `,
	` feat `}
