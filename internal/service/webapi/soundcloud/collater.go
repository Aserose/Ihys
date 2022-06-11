package soundcloud

import (
	"IhysBestowal/internal/datastruct"
	"sync"
)

const (
	soundcloudTitle = "soundcloud"
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
			quantityFlow:            1,
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

func (c collater) getSimilarParallel(sourceData datastruct.AudioItems) datastruct.AudioItems {
	if c.options.maxAudioAmountPerSource <= 0 {
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

	for sourceIndexFrom := 0; sourceIndexFrom <= len(sourceData.Items); sourceIndexFrom += c.options.quantityFlow {
		if len(sourceData.Items[sourceIndexFrom:]) <= c.options.quantityFlow {
			ch <- c.getSimilar(sourceData.Items[sourceIndexFrom:])
			break
		}

		wg.Add(1)
		go func(source []datastruct.AudioItem) {
			defer wg.Done()
			ch <- c.getSimilar(source)
		}(sourceData.Items[sourceIndexFrom : sourceIndexFrom+c.options.quantityFlow])
	}
	wg.Wait()
	close(ch)
	closed <- true
	close(closed)

	return datastruct.AudioItems{
		Items: collectedSimilar,
		From:  soundcloudTitle,
	}
}

func (c collater) getSimilar(sourceItems []datastruct.AudioItem) (result []datastruct.AudioItem) {
	for _, items := range sourceItems {
		result = append(result, c.collate(c.parser.getSimilar(items.Artist, items.Title))...)
	}
	return
}

func (c collater) collateWithoutArtistStrain(data []datastruct.AudioItem) []datastruct.AudioItem {
	if len(data) >= c.options.maxAudioAmountPerSource {
		return data[:c.options.maxAudioAmountPerSource]
	}
	return data
}

func (c collater) collateWithArtistStrain(data []datastruct.AudioItem) []datastruct.AudioItem {
	numberOfArtistSongs := make(map[string]int)
	var artistName string

	for i := 0; i < len(data)-1; i++ {
		artistName = data[i].Artist

		if artistName == data[i+1].Artist {
			numberOfArtistSongs[artistName]++

			if numberOfArtistSongs[artistName] >= c.options.maxAudioAmountPerArtist {
				data = append(data[:i], data[i+1:]...)
				i--
			}
		}
	}

	if len(data) >= c.options.maxAudioAmountPerSource {
		return data[:c.options.maxAudioAmountPerSource]
	}
	return data
}
