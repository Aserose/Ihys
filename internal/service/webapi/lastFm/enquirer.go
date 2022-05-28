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

type iEnquirer interface {
	getSimilarTracks(map[string]interface{}) lastfm.TrackGetSimilar
	getTopTracks(artistNames []string, numberOfTracksPerArtist int) datastruct.AudioItems
	getSimilarArtists(artistName string, limit int) []string
}

type enquirer struct {
	api *lastfm.Api
	log customLogger.Logger
}

func newEnquirer(log customLogger.Logger, cfg config.LastFM, repo repository.Repository) iEnquirer {
	return &enquirer{
		api: lastfm.New(cfg.Key, cfg.Secret),
		log: log,
	}
}

func (l enquirer) getSimilarTracks(queryParams map[string]interface{}) lastfm.TrackGetSimilar {
	similiar, _ := l.api.Track.GetSimilar(queryParams)

	return similiar
}

func (l enquirer) getTopTracks(artistNames []string, numberOfTracksPerArtist int) datastruct.AudioItems {
	wg := &sync.WaitGroup{}
	res := datastruct.AudioItems{}

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
				if i >= numberOfTracksPerArtist {
					break
				}

				res.Items = append(res.Items, datastruct.AudioItem{
					Artist: artistName,
					Title:  track.Name,
				})
			}
		}(artistName)
	}
	wg.Wait()

	res.From = `\top`

	return res
}

func (l enquirer) getSimilarArtists(artistName string, limit int) []string {
	if limit <= 0 {
		return []string{}
	}

	var (
		result []string
		wg     = &sync.WaitGroup{}
	)

	resp := lastfm.ArtistGetSimilar{}

	request := func(artistName string) []string {
		resp, _ = l.api.Artist.GetSimilar(map[string]interface{}{
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
						result = append(result, request(name)...)
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
		result = append(result, request(artistName)...)
	}

	if result == nil {
		return []string{}
	}

	return result
}

var enumTypes = []string{
	`, `,
	` feat. `,
	` feat `}
