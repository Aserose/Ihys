package yaMusic

import (
	"IhysBestowal/internal/datastruct"
	"sync"
)

type collater struct {
	parser
	options
	collate func(data []datastruct.AudioItem) []datastruct.AudioItem
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
	if m.options.maxAudioAmountPerSource <= 0 {
		return datastruct.AudioItems{}
	}

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
		if len(sourceData.Items[sourceIndexFrom:]) <= m.options.quantityFlow {
			ch <- m.getSimilar(sourceData.Items[sourceIndexFrom:])
			break
		}

		wg.Add(1)
		go func(source []datastruct.AudioItem) {
			defer wg.Done()
			ch <- m.getSimilar(source)
		}(sourceData.Items[sourceIndexFrom : sourceIndexFrom+m.options.quantityFlow])
	}
	wg.Wait()
	close(ch)
	closed <- true
	close(closed)

	return datastruct.AudioItems{
		Items: collectedSimilar,
		From:  "yaMusic",
	}
}

func (m collater) getSimilar(sourceItems []datastruct.AudioItem) (result []datastruct.AudioItem) {
	for _, item := range sourceItems {
		result = append(result, m.collate(m.parser.getSimilar(item.Artist, item.Title))...)
	}
	return
}

func (m collater) collateWithoutArtistStrain(data []datastruct.AudioItem) []datastruct.AudioItem {
	if len(data) >= m.options.maxAudioAmountPerSource {
		return data[:m.options.maxAudioAmountPerSource]
	}
	return data
}

func (m collater) collateWithArtistStrain(data []datastruct.AudioItem) []datastruct.AudioItem {
	numberOfArtistSongs := make(map[string]int)
	var artistName string

	for i := 0; i < len(data)-1; i++ {
		artistName = data[i].Artist

		if artistName == data[i+1].Artist {
			numberOfArtistSongs[artistName]++

			if numberOfArtistSongs[artistName] >= m.options.maxAudioAmountPerArtist {
				data = append(data[:i], data[i+1:]...)
				i--
			}
		}
	}

	if len(data) >= m.options.maxAudioAmountPerSource {
		return data[:m.options.maxAudioAmountPerSource]
	}
	return data
}
