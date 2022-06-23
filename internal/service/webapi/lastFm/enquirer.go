package lastFm

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/repository"
	"IhysBestowal/pkg/customLogger"
	"github.com/goccy/go-json"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
	"sync"
)

type enquirer struct {
	apiKey     string
	httpClient *fasthttp.Client
	log        customLogger.Logger
}

func newEnquirer(log customLogger.Logger, cfg config.LastFM, repo repository.Repository) enquirer {
	return enquirer{
		apiKey:     cfg.Key,
		httpClient: &fasthttp.Client{},
		log:        log,
	}
}

func (l enquirer) sendRequest(req *fasthttp.Request) []byte {
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := l.httpClient.Do(req, resp)
	if err != nil {
		l.log.Warn(l.log.CallInfoStr(), err.Error())
	}

	return resp.Body()
}

func (l enquirer) getAudio(query string) datastruct.AudioItem {
	if query == "" {
		return datastruct.AudioItem{}
	}

	resp := datastruct.LastFMSearchTrackResult{}

	uri := fasthttp.AcquireURI()
	uri.Parse(nil, []byte(baseUrl))
	uri.QueryArgs().Add("method", searchTrack)
	uri.QueryArgs().Add("limit", "1")
	uri.QueryArgs().Add("track", query)
	uri.QueryArgs().Add("api_key", l.apiKey)
	uri.QueryArgs().Add("format", jsonFrmt)

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(fasthttp.MethodGet)
	req.SetURI(uri)
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseURI(uri)
	}()
	json.Unmarshal(l.sendRequest(req), &resp)

	return datastruct.AudioItem{
		Artist: resp.Results.TrackMatches.Tracks[0].Artist,
		Title:  resp.Results.TrackMatches.Tracks[0].Name,
	}
}

func (l enquirer) getSimilarTracks(artist, track string) datastruct.AudioItems {
	resp := datastruct.LastFMUnmr{}

	uri := fasthttp.AcquireURI()
	uri.Parse(nil, []byte(baseUrl))
	uri.QueryArgs().Add("method", getSimilarTrack)
	uri.QueryArgs().Add("artist", artist)
	uri.QueryArgs().Add("track", track)
	uri.QueryArgs().Add("api_key", l.apiKey)
	uri.QueryArgs().Add("format", jsonFrmt)

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(fasthttp.MethodGet)
	req.SetURI(uri)
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseURI(uri)
	}()

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

	resp := datastruct.LastFMUnmr{}
	trackList := make([]datastruct.AudioItem, len(artistNames)*numberOfTracksPerArtist)
	wg := &sync.WaitGroup{}
	ch := make(chan datastruct.AudioItem)
	closed := make(chan bool)

	request := func(artistName string) datastruct.LastFMTopTracks {
		uri := fasthttp.AcquireURI()
		uri.Parse(nil, []byte(baseUrl))
		uri.QueryArgs().Add("method", getTopTrack)
		uri.QueryArgs().Add("artist", artistName)
		uri.QueryArgs().Add("api_key", l.apiKey)
		uri.QueryArgs().Add("format", jsonFrmt)

		req := fasthttp.AcquireRequest()
		req.Header.SetMethod(fasthttp.MethodGet)
		req.SetURI(uri)
		defer func() {
			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseURI(uri)
		}()

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

	for _, artistName := range artistNames {
		wg.Add(1)
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

	resp := datastruct.LastFMUnmr{}
	artistList := []string{}
	wg := &sync.WaitGroup{}
	ch := make(chan []string)
	closed := make(chan bool)

	request := func(artistName string) []string {
		uri := fasthttp.AcquireURI()
		uri.Parse(nil, []byte(baseUrl))
		uri.QueryArgs().Add("method", getSimilarArtist)
		uri.QueryArgs().Add("limit", strconv.Itoa(limit))
		uri.QueryArgs().Add("artist", artistName)
		uri.QueryArgs().Add("api_key", l.apiKey)
		uri.QueryArgs().Add("format", jsonFrmt)
		uri.QueryArgs().Add("autocorrect", "1")

		req := fasthttp.AcquireRequest()
		req.Header.SetMethod(fasthttp.MethodGet)
		req.SetURI(uri)
		defer func() {
			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseURI(uri)
		}()

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
				artistList = append(artistList, names...)
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

	if artistList == nil {
		return []string{}
	}

	return artistList
}

var enumTypes = []string{
	`, `,
	` feat. `,
	` feat `}
