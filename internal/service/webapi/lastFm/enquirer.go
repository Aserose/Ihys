package lastFm

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/repository"
	"IhysBestowal/pkg/customLogger"
	"github.com/goccy/go-json"
	"github.com/valyala/fasthttp"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

type enquirer struct {
	apiKey      string
	sendRequest func(req *fasthttp.Request) []byte
	log         customLogger.Logger
}

func newEnquirer(log customLogger.Logger, cfg config.LastFM, repo repository.Repository) enquirer {
	httpClient := &fasthttp.Client{}

	sendRequest := func(req *fasthttp.Request) []byte {
		b := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(b)

		err := httpClient.Do(req, b)
		if err != nil {
			log.Warn(log.CallInfoStr(), err.Error())
		}

		return b.Body()
	}

	return enquirer{
		apiKey:      cfg.Key,
		sendRequest: sendRequest,
		log:         log,
	}
}

func (l enquirer) getSimilarTracks(artist, track string) datastruct.LastFMSimilarTracks {
	resp := datastruct.LastFMUnmr{}

	base, _ := url.Parse(baseUri)
	values := url.Values{}
	values.Add("method", getSimilarTrack)
	values.Add("artist", artist)
	values.Add("track", track)
	values.Add("api_key", l.apiKey)
	values.Add("format", jsonFrmt)
	base.RawQuery = values.Encode()

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(fasthttp.MethodGet)
	req.SetRequestURI(base.String())
	defer fasthttp.ReleaseRequest(req)

	json.Unmarshal(l.sendRequest(req), &resp)

	return resp.LastFMSimilarTracks
}

func (l enquirer) getTopTracks(artistNames []string, numberOfTracksPerArtist int) datastruct.AudioItems {
	if artistNames == nil || numberOfTracksPerArtist <= 0 {
		return datastruct.AudioItems{}
	}

	resp := datastruct.LastFMUnmr{}
	wg := &sync.WaitGroup{}
	res := make([]datastruct.AudioItem, len(artistNames)*numberOfTracksPerArtist)
	ch := make(chan datastruct.AudioItem)
	closed := make(chan bool)

	go func() {
		j := 0
		for {
			select {
			case inc, ok := <-ch:
				if !ok {
					continue
				}

				res[j] = inc
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

			base, _ := url.Parse(baseUri)
			values := url.Values{}
			values.Add("method", getTopTrack)
			values.Add("artist", artistName)
			values.Add("api_key", l.apiKey)
			values.Add("format", jsonFrmt)
			base.RawQuery = values.Encode()

			req := fasthttp.AcquireRequest()
			req.Header.SetMethod(fasthttp.MethodGet)
			req.SetRequestURI(base.String())
			defer fasthttp.ReleaseRequest(req)

			json.Unmarshal(l.sendRequest(req), &resp)

			for i, track := range resp.LastFMTopTracks.Tracks {
				t := track
				if i >= numberOfTracksPerArtist-1 {
					break
				}
				ch <- datastruct.AudioItem{
					Artist: artistName,
					Title:  t.Name,
				}
			}
		}(artistName)
	}
	wg.Wait()
	close(ch)
	closed <- true

	return datastruct.AudioItems{
		Items: res,
		From:  `\top`,
	}
}

func (l enquirer) getSimilarArtists(artistName string, limit int) []string {
	if limit <= 0 {
		return []string{}
	}

	resp := datastruct.LastFMUnmr{}
	res := []string{}
	wg := &sync.WaitGroup{}
	ch := make(chan []string)
	closed := make(chan bool)

	go func() {
		for {
			select {
			case inc, ok := <-ch:
				if !ok {
					continue
				}
				res = append(res, inc...)
			case <-closed:
				return
			}
		}
	}()

	request := func(artistName string) []string {

		base, _ := url.Parse(baseUri)
		values := url.Values{}
		values.Add("method", getSimilarArtist)
		values.Add("limit", strconv.Itoa(limit))
		values.Add("artist", artistName)
		values.Add("api_key", l.apiKey)
		values.Add("format", jsonFrmt)
		values.Add("autocorrect", "1")
		base.RawQuery = values.Encode()

		req := fasthttp.AcquireRequest()
		req.Header.SetMethod(fasthttp.MethodGet)
		req.SetRequestURI(base.String())
		defer fasthttp.ReleaseRequest(req)

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

	if res == nil {
		return []string{}
	}

	return res
}

var enumTypes = []string{
	`, `,
	` feat. `,
	` feat `}
