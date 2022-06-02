package lastFm

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/repository"
	"IhysBestowal/pkg/customLogger"
	"github.com/shkh/lastfm-go/lastfm"
	"strings"
	"sync"
)

type enquirer struct {
	api *lastfm.Api
	log customLogger.Logger
}

func newEnquirer(log customLogger.Logger, cfg config.LastFM, repo repository.Repository) enquirer {
	return enquirer{
		api: lastfm.New(cfg.Key, cfg.Secret),
		log: log,
	}
}

func (l enquirer) getSimilarTracks(queryParams map[string]interface{}) lastfm.TrackGetSimilar {
	similiar, _ := l.api.Track.GetSimilar(queryParams)

	return similiar
}

func (l enquirer) getTopTracks(artistNames []string, numberOfTracksPerArtist int) datastruct.AudioItems {
	if artistNames == nil || numberOfTracksPerArtist <= 0 {
		return datastruct.AudioItems{}
	}

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

			tracks, err := l.api.Artist.GetTopTracks(map[string]interface{}{
				"artist": artistName,
			})
			if err != nil {
				l.log.Error(l.log.CallInfoStr(), err.Error())
			}

			for i, track := range tracks.Tracks {
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
		resp, _ := l.api.Artist.GetSimilar(map[string]interface{}{
			"limit":       limit,
			"artist":      artistName,
			"autocorrect": 1,
		})

		if resp.Similars == nil {
			return []string{}
		}
		artistList := make([]string, len(resp.Similars))
		for i, r := range resp.Similars {
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
