package yaMusic

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
	"strings"
	"sync"
)

type iEnquirer interface {
	getSimiliar(sourceAudios datastruct.AudioItems) datastruct.AudioItems
	getAudio(query string) datastruct.AudioItem
}

type enquirer struct {
	iParser
	options
	collate func(sourceAudio datastruct.AudioItem) []datastruct.AudioItem
}

func newEnquirer(log customLogger.Logger, opts ...processingOptions) iEnquirer {
	enq := enquirer{
		iParser: newParser(log),
		options: options{
			quantityFlow:         3,
			numberOfSong:         2,
			audioAmountPerSource: 3,
			audioAmountPerArtist: 1,
		},
	}

	for _, opt := range opts {
		opt(&enq.options)
	}

	if enq.audioAmountPerArtist == 0 {
		enq.collate = enq.collateWithoutArtistStrain
	} else {
		enq.collate = enq.collateWithArtistStrain
	}

	return enq
}

func (m enquirer) getAudio(query string) datastruct.AudioItem {
	return m.iParser.getAudio(query)
}

func (m enquirer) getSimiliar(sourceData datastruct.AudioItems) (result datastruct.AudioItems) {
	wg := sync.WaitGroup{}

	for sourceDataFrom := 0; sourceDataFrom <= len(sourceData.Items); sourceDataFrom += m.options.quantityFlow {
		var items []datastruct.AudioItem

		if len(sourceData.Items[sourceDataFrom:]) < m.options.quantityFlow {
			items = sourceData.Items[sourceDataFrom:]
			m.discover(items)
			break
		}

		items = sourceData.Items[sourceDataFrom : sourceDataFrom+m.options.quantityFlow]

		wg.Add(1)
		go func() {
			defer wg.Done()
			result.Items = append(result.Items, m.discover(items)...)
		}()
	}
	wg.Wait()
	result.From = "yaMusic"

	return
}

func (m enquirer) discover(sourceItems []datastruct.AudioItem) (items []datastruct.AudioItem) {
	for _, sourceAudio := range sourceItems {
		items = append(items, m.collate(sourceAudio)...)
	}

	return
}

func (m enquirer) collateWithoutArtistStrain(sourceAudio datastruct.AudioItem) (items []datastruct.AudioItem) {
	addToResultItems := func(item datastruct.AudioItem) {
		items = append(items, item)
	}

	for j, sim := range m.iParser.getDecodedData(sourceAudio.Artist, sourceAudio.Title).Sidebar.SimilarTracks {
		if j > m.options.audioAmountPerSource {
			break
		}
		s := sim

		addToResultItems(datastruct.AudioItem{
			Artist: m.writeArtistName(s.Artists),
			Title:  s.Title,
		})
	}

	return
}

func (m enquirer) collateWithArtistStrain(sourceAudio datastruct.AudioItem) (items []datastruct.AudioItem) {
	isTheArtistOnTheResult := func(artist string) bool {
		numberOfArtistSongs := map[string]int{}

		for _, item := range items {
			if item.Artist == artist {
				numberOfArtistSongs[artist]++
				if numberOfArtistSongs[artist] > m.options.audioAmountPerArtist {
					return true
				}
			}
		}
		return false
	}

	addToResultItems := func(item datastruct.AudioItem) (artistAlreadyOnTheList bool) {
		if isTheArtistOnTheResult(item.Artist) {
			return true
		}
		items = append(items, item)
		return false
	}

	j := 0
	for _, sim := range m.iParser.getDecodedData(sourceAudio.Artist, sourceAudio.Title).Sidebar.SimilarTracks {
		if j > m.options.audioAmountPerSource {
			break
		}
		s := sim

		artistAlreadyOnTheList := addToResultItems(datastruct.AudioItem{
			Artist: m.writeArtistName(s.Artists),
			Title:  s.Title,
		})
		if artistAlreadyOnTheList {
			j--
		}
		j++
	}

	return
}

func (m enquirer) writeArtistName(artists []struct {
	Name string `json:"name"`
}) string {
	artistName := strings.Builder{}

	if len(artists) > 1 {
		for i, artist := range artists {
			artistName.WriteString(artist.Name)
			if i < len(artists)-1 {
				artistName.WriteString(", ")
			}
		}
	} else {
		artistName.WriteString(artists[0].Name)
	}

	return artistName.String()
}
