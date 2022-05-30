package yaMusic

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
	"strings"
	"sync"
)

type iCollater interface {
	getSimiliar(sourceAudios datastruct.AudioItems) datastruct.AudioItems
	getAudio(query string) datastruct.AudioItem
}

type collater struct {
	iParser
	options
	collate func(sourceAudio datastruct.AudioItem) []datastruct.AudioItem
}

func newCollater(log customLogger.Logger, opts ...processingOptions) iCollater {
	enq := collater{
		iParser: newParser(log),
		options: options{
			quantityFlow:            3,
			maxAudioAmountPerSource: 3,
			maxAudioAmountPerArtist: 1,
		},
	}

	for _, opt := range opts {
		opt(&enq.options)
	}

	if enq.maxAudioAmountPerArtist == 0 {
		enq.collate = enq.collateWithoutArtistStrain
	} else {
		enq.collate = enq.collateWithArtistStrain
	}

	return enq
}

func (m collater) getAudio(query string) datastruct.AudioItem {
	return m.iParser.getAudio(query)
}

func (m collater) getSimiliar(sourceData datastruct.AudioItems) (result datastruct.AudioItems) {
	wg := sync.WaitGroup{}

	for sourceIndexFrom := 0; sourceIndexFrom <= len(sourceData.Items); sourceIndexFrom += m.options.quantityFlow {
		var sourceAudio []datastruct.AudioItem

		if len(sourceData.Items[sourceIndexFrom:]) < m.options.quantityFlow {
			sourceAudio = sourceData.Items[sourceIndexFrom:]
			result.Items = append(result.Items, m.collectSimiliar(sourceAudio)...)
			break
		}

		sourceAudio = sourceData.Items[sourceIndexFrom : sourceIndexFrom+m.options.quantityFlow]

		wg.Add(1)
		go func() {
			defer wg.Done()
			result.Items = append(result.Items, m.collectSimiliar(sourceAudio)...)
		}()

	}

	wg.Wait()

	result.From = "yaMusic"

	return
}

func (m collater) collectSimiliar(sourceItems []datastruct.AudioItem) (result []datastruct.AudioItem) {
	for _, item := range sourceItems {
		result = append(result, m.collate(item)...)
	}

	return
}

func (m collater) collateWithoutArtistStrain(sourceAudio datastruct.AudioItem) (result []datastruct.AudioItem) {
	for j, sim := range m.iParser.getSimiliars(sourceAudio.Artist, sourceAudio.Title).YaMSidebar.SimilarTracks {
		if j >= m.options.maxAudioAmountPerSource { break }
		s := sim

		result = append(result, datastruct.AudioItem{
			Artist: m.writeArtistName(s.Artists),
			Title:  s.Title,
		})
	}

	return
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

	j := 0
	for _, sim := range m.iParser.getSimiliars(sourceAudio.Artist, sourceAudio.Title).YaMSidebar.SimilarTracks {
		if j >= m.options.maxAudioAmountPerSource { break }
		s := sim

		limitReached := addToResultItems(datastruct.AudioItem{
			Artist: m.writeArtistName(s.Artists),
			Title:  s.Title,
		}); if !limitReached { j++ }
	}

	return
}

func (m collater) writeArtistName(artists []datastruct.YaMArtists) string {
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
