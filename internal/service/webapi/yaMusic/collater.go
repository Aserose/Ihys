package yaMusic

import (
	"IhysBestowal/internal/datastruct"
	"sync"
)

type collater struct {
	parser
	options
	collate func(sourceAudio datastruct.AudioItem) []datastruct.AudioItem
}

func newCollater(p parser, opts ...ProcessingOptions) collater {
	cltr := collater{
		parser: p,
		options: options{
			quantityFlow:            3,
			maxAudioAmountPerSource: 3,
			maxAudioAmountPerArtist: 1,
		},
	}

	if opts != nil {
		for _, opt := range opts {
			opt(&cltr.options)
		}
	}

	if cltr.options.maxAudioAmountPerArtist == 0 {
		cltr.collate = cltr.collateWithoutArtistStrain
	} else {
		cltr.collate = cltr.collateWithArtistStrain
	}

	return cltr
}

func (m collater) getSimilarParallel(sourceData datastruct.AudioItems) datastruct.AudioItems {
	wg := &sync.WaitGroup{}
	collectedSimilar := []datastruct.AudioItem{}
	ch := make(chan []datastruct.AudioItem)
	closed := make(chan bool)

	go func() {
		for {
			select {
			case sim, ok := <-ch:
				if !ok {
					continue
				}
				collectedSimilar = append(collectedSimilar, sim...)
			case <-closed:
				return
			}
		}
	}()

	for sourceIndexFrom := 0; sourceIndexFrom <= len(sourceData.Items); sourceIndexFrom += m.options.quantityFlow {
		var sourceAudio []datastruct.AudioItem

		if len(sourceData.Items[sourceIndexFrom:]) < m.options.quantityFlow {
			sourceAudio = sourceData.Items[sourceIndexFrom:]
			ch <- m.getSimilar(sourceAudio)
			break
		}

		sourceAudio = sourceData.Items[sourceIndexFrom : sourceIndexFrom+m.options.quantityFlow]

		wg.Add(1)
		go func() {
			defer wg.Done()
			ch <- m.getSimilar(sourceAudio)
		}()
	}
	wg.Wait()
	close(ch)
	closed <- true

	return datastruct.AudioItems{
		Items: collectedSimilar,
		From:  "yaMusic",
	}
}

func (m collater) getSimilar(sourceItems []datastruct.AudioItem) (result []datastruct.AudioItem) {
	for _, item := range sourceItems {
		result = append(result, m.collate(item)...)
	}

	return
}

func (m collater) collateWithoutArtistStrain(sourceAudio datastruct.AudioItem) (items []datastruct.AudioItem) {
	simTracks := m.parser.getSimilar(sourceAudio.Artist, sourceAudio.Title).YaMSidebar.SimilarTracks
	items = make([]datastruct.AudioItem, len(simTracks))

	for j, sim := range simTracks {
		if j >= m.options.maxAudioAmountPerSource {
			break
		}
		s := sim

		items[j] = datastruct.AudioItem{
			Artist: m.writeArtistName(s.Artists),
			Title:  s.Title,
		}
	}

	return items
}

func (m collater) collateWithArtistStrain(sourceAudio datastruct.AudioItem) (items []datastruct.AudioItem) {
	artistSongLimitReached := func(artist string) bool {
		numberOfArtistSongs := map[string]int{}

		for _, item := range items {
			if item.Artist == artist {
				numberOfArtistSongs[artist]++
				if numberOfArtistSongs[artist] > m.options.maxAudioAmountPerArtist {
					return true
				}
			}
		}
		return false
	}

	addToResultItems := func(item datastruct.AudioItem) (limitReached bool) {
		if artistSongLimitReached(item.Artist) {
			return true
		}

		items = append(items, item)
		return false
	}

	i := 0
	for _, sim := range m.parser.getSimilar(sourceAudio.Artist, sourceAudio.Title).YaMSidebar.SimilarTracks {
		if i >= m.options.maxAudioAmountPerSource {
			break
		}
		s := sim

		limitReached := addToResultItems(datastruct.AudioItem{
			Artist: m.writeArtistName(s.Artists),
			Title:  s.Title,
		})
		if !limitReached {
			i++
		}
	}

	return
}

func (m collater) writeArtistName(artists []datastruct.YaMArtists) (result string) {
	if len(artists) > 1 {
		for i, artist := range artists {
			result += artist.Name
			if i < len(artists)-1 {
				result += ", "
			}
		}
	} else {
		result += artists[0].Name
	}

	return
}
