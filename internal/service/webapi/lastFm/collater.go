package lastFm

import (
	"IhysBestowal/internal/datastruct"
	"github.com/shkh/lastfm-go/lastfm"
	"sync"
)

type iCollater interface {
	getSimiliars(userId int64, sourceData datastruct.AudioItems) (toResult datastruct.AudioItems)
}

type collater struct {
	options
	collate func(data interface{}) []datastruct.AudioItem
	iEnquirer
}

func newCollater(enq iEnquirer, opts ...processingOptions) iCollater {
	cl := collater{
		options: options{
			quantityFlow:            3,
			maxAudioAmountPerSource: 3,
			maxAudioAmountPerArtist: 1,
			maxNumSimiliarArtists:   35,
			maxNumTopPerArtist:      10,
		},
		iEnquirer: enq,
	}

	for _, o := range opts {
		o(&cl.options)
	}

	if cl.options.maxAudioAmountPerArtist == 0 {
		cl.collate = cl.collateWithoutArtistStrain
	} else {
		cl.collate = cl.collateWithArtistStrain
	}

	return cl
}

func (c collater) getSimiliars(userId int64, sourceData datastruct.AudioItems) (result datastruct.AudioItems) {
	wg := sync.WaitGroup{}

	for sourceDataFrom := 0; sourceDataFrom <= len(sourceData.Items); sourceDataFrom += c.options.quantityFlow {
		var audioItems []datastruct.AudioItem

		if len(sourceData.Items[sourceDataFrom:]) < c.options.quantityFlow {
			audioItems = sourceData.Items[sourceDataFrom:]
			result.Items = append(result.Items, c.collectSimiliar(audioItems)...)
			break
		}

		audioItems = sourceData.Items[sourceDataFrom : sourceDataFrom+c.options.quantityFlow]

		wg.Add(1)
		go func() {
			defer wg.Done()
			result.Items = append(result.Items, c.collectSimiliar(audioItems)...)
		}()
	}
	wg.Wait()

	result.From = "lastFm"

	return
}

func (c collater) collectSimiliar(sourceItems []datastruct.AudioItem) (resultItems []datastruct.AudioItem) {
	queryParams := make(map[string]interface{})

	for _, d := range sourceItems {
		queryParams["artist"], queryParams["track"] = d.Artist, d.Title

		similiar := c.iEnquirer.getSimilarTracks(queryParams)

		switch similiar.Tracks != nil {
		case true:
			sim := c.collate(similiar)
			resultItems = append(resultItems, sim...)

			if len(sim) < c.maxAudioAmountPerSource {
				resultItems = append(resultItems, c.collate(c.iEnquirer.getTopTracks(
							c.iEnquirer.getSimilarArtists(d.Artist, c.maxNumSimiliarArtists),
							c.maxAudioAmountPerSource-len(sim)))...)
			}
		case false:
			resultItems = append(resultItems, c.collate(c.iEnquirer.getTopTracks(
						c.iEnquirer.getSimilarArtists(d.Artist, c.maxNumSimiliarArtists),
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
			if i >= c.maxAudioAmountPerSource { break }

			addToResultItems(s.Artist.Name, s.Name)
		}
	case []datastruct.AudioItems:
		for i, s := range data.([]datastruct.LastFMResponse) {
			if i >= c.maxAudioAmountPerSource { break }

			addToResultItems(s.Artist, s.Title)
		}
	}

	return audioItems
}

func (c collater) collateWithArtistStrain(data interface{}) []datastruct.AudioItem {
	audioItems := []datastruct.AudioItem{}

	isTheArtistOnTheResult := func(artist string) bool {
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

	addToResultItems := func(artist, title string) bool {
		if isTheArtistOnTheResult(artist) {
			return false
		}

		audioItems = append(audioItems, datastruct.AudioItem{
			Artist: artist,
			Title:  title,
		})
		return true
	}

	switch data.(type) {
	case lastfm.TrackGetSimilar:
		i := 0
		for _, s := range data.(lastfm.TrackGetSimilar).Tracks {
			if i >= c.options.maxAudioAmountPerSource { break }

			if !addToResultItems(s.Artist.Name, s.Name) { i-- }
			i++
		}
	case datastruct.AudioItems:
		i := 0
		for _, s := range data.(datastruct.AudioItems).Items {
			if i >= c.options.maxAudioAmountPerSource { break }

			if !addToResultItems(s.Artist, s.Title) { i-- }
			i++
		}
	}

	return audioItems
}
