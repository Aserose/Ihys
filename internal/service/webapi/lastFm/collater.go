package lastFm

import (
	"IhysBestowal/internal/datastruct"
	"github.com/shkh/lastfm-go/lastfm"
	"sync"
)

type collater struct {
	options
	collate func(data interface{}) []datastruct.AudioItem
	enquirer
}

func newCollater(enq enquirer, opts ...ProcessingOptions) collater {
	cl := collater{
		options: options{
			quantityFlow:            3,
			maxAudioAmountPerSource: 3,
			maxAudioAmountPerArtist: 1,
			maxNumSimiliarArtists:   35,
			maxNumTopPerArtist:      10,
		},
		enquirer: enq,
	}

	if opts != nil {
		for _, o := range opts {
			o(&cl.options)
		}
	}

	if cl.options.maxAudioAmountPerArtist == 0 {
		cl.collate = cl.collateWithoutArtistStrain
	} else {
		cl.collate = cl.collateWithArtistStrain
	}

	return cl
}

func (c collater) getSimilarParallel(userId int64, sourceData datastruct.AudioItems) datastruct.AudioItems {
	wg := &sync.WaitGroup{}
	res := []datastruct.AudioItem{}
	ch := make(chan []datastruct.AudioItem)
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

	for sourceDataFrom := 0; sourceDataFrom <= len(sourceData.Items); sourceDataFrom += c.options.quantityFlow {
		var audioItems []datastruct.AudioItem

		if len(sourceData.Items[sourceDataFrom:]) < c.options.quantityFlow {
			audioItems = sourceData.Items[sourceDataFrom:]
			ch <- c.getSimilar(audioItems)
			break
		}

		audioItems = sourceData.Items[sourceDataFrom : sourceDataFrom+c.options.quantityFlow]

		wg.Add(1)
		go func() {
			defer wg.Done()
			ch <- c.getSimilar(audioItems)
		}()
	}

	wg.Wait()
	close(ch)
	closed <- true

	return datastruct.AudioItems{
		Items: res,
		From:  "lastFm",
	}
}

func (c collater) getSimilar(sourceItems []datastruct.AudioItem) (resultItems []datastruct.AudioItem) {
	queryParams := make(map[string]interface{})

	for _, d := range sourceItems {
		queryParams["artist"], queryParams["track"] = d.Artist, d.Title

		similiar := c.enquirer.getSimilarTracks(queryParams)

		switch similiar.Tracks != nil {
		case true:
			sim := c.collate(similiar)
			resultItems = append(resultItems, sim...)

			if len(sim) < c.maxAudioAmountPerSource {
				resultItems = append(resultItems, c.collate(c.enquirer.getTopTracks(
					c.enquirer.getSimilarArtists(d.Artist, c.maxNumSimiliarArtists),
					c.maxAudioAmountPerSource-len(sim)))...)
			}
		case false:
			resultItems = append(resultItems, c.collate(c.enquirer.getTopTracks(
				c.enquirer.getSimilarArtists(d.Artist, c.maxNumSimiliarArtists),
				c.maxNumTopPerArtist))...)
		}
	}

	return
}

func (c collater) collateWithoutArtistStrain(data interface{}) []datastruct.AudioItem {
	audioItems := []datastruct.AudioItem{}

	addToResultItems := func(artist, title string) {
		audioItems = append(audioItems, datastruct.AudioItem{
			Artist: artist,
			Title:  title,
		})
	}

	switch data.(type) {
	case lastfm.TrackGetSimilar:
		for i, s := range data.(lastfm.TrackGetSimilar).Tracks {
			if i >= c.maxAudioAmountPerSource {
				break
			}

			addToResultItems(s.Artist.Name, s.Name)
		}
	case []datastruct.AudioItems:
		for i, s := range data.([]datastruct.LastFMResponse) {
			if i >= c.maxAudioAmountPerSource {
				break
			}

			addToResultItems(s.Artist, s.Title)
		}
	}

	return audioItems
}

func (c collater) collateWithArtistStrain(data interface{}) []datastruct.AudioItem {
	audioItems := []datastruct.AudioItem{}

	artistSongLimitReached := func(artist string) bool {
		numberOfArtistSongs := map[string]int{}

		for _, item := range audioItems {
			if item.Artist == artist {
				numberOfArtistSongs[artist]++
				if numberOfArtistSongs[artist] >= c.options.maxAudioAmountPerArtist {
					return true
				}
			}
		}
		return false
	}

	addToResultItems := func(artist, title string) (limitReached bool) {
		if artistSongLimitReached(artist) {
			return true
		}

		audioItems = append(audioItems, datastruct.AudioItem{
			Artist: artist,
			Title:  title,
		})
		return false
	}

	switch data.(type) {
	case lastfm.TrackGetSimilar:
		i := 0
		for _, s := range data.(lastfm.TrackGetSimilar).Tracks {
			if i >= c.options.maxAudioAmountPerSource {
				break
			}

			limitReached := addToResultItems(s.Artist.Name, s.Name)
			if !limitReached {
				i++
			}

		}
	case datastruct.AudioItems:
		i := 0
		for _, s := range data.(datastruct.AudioItems).Items {
			if i >= c.options.maxAudioAmountPerSource {
				break
			}

			limitReached := addToResultItems(s.Artist, s.Title)
			if !limitReached {
				i++
			}
		}
	}

	return audioItems
}
