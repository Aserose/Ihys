package lastFm

import (
	"IhysBestowal/internal/datastruct"
	"sort"
	"sync"
)

type collater struct {
	options
	collate func(data datastruct.AudioItems) []datastruct.AudioItem
	enquirer
}

func newCollater(enq enquirer, opts ...ProcessingOptions) collater {
	cl := collater{
		options: options{
			quantityFlow:            3,
			maxAudioAmountPerSource: 3,
			maxAudioAmountPerArtist: 1,
			maxNumSimilarArtists:    35,
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
	if c.options.maxAudioAmountPerSource <= 0 || c.options.maxAudioAmountPerArtist <= 0 {
		return datastruct.AudioItems{}
	}

	res := []datastruct.AudioItem{}
	wg := &sync.WaitGroup{}
	ch := make(chan []datastruct.AudioItem)
	closed := make(chan bool)

	go func() {
		for {
			select {
			case similar, ok := <-ch:
				if !ok {
					continue
				}
				res = append(res, similar...)
			case <-closed:
				return
			}
		}
	}()

	for sourceDataFrom := 0; sourceDataFrom <= len(sourceData.Items); sourceDataFrom += c.options.quantityFlow {
		if len(sourceData.Items[sourceDataFrom:]) <= c.options.quantityFlow {
			ch <- c.getSimilar(sourceData.Items[sourceDataFrom:])
			break
		}

		wg.Add(1)
		go func(source []datastruct.AudioItem) {
			defer wg.Done()
			ch <- c.getSimilar(source)
		}(sourceData.Items[sourceDataFrom : sourceDataFrom+c.options.quantityFlow])
	}
	wg.Wait()
	close(ch)
	closed <- true
	close(closed)

	return datastruct.AudioItems{
		Items: res,
		From:  SourceFrom,
	}
}

func (c collater) getSimilar(sourceItems []datastruct.AudioItem) (resultItems []datastruct.AudioItem) {
	for _, d := range sourceItems {
		similar := c.enquirer.getSimilarTracks(d.Artist, d.Title)

		switch similar.Items != nil {
		case true:
			collatedSimilar := c.collate(similar)
			resultItems = append(resultItems, collatedSimilar...)

			if len(collatedSimilar) < c.maxAudioAmountPerSource {
				resultItems = append(resultItems, c.collate(
					c.enquirer.getTopTracks(
						c.enquirer.getSimilarArtists(d.Artist, c.maxAudioAmountPerSource-len(collatedSimilar)),
						c.maxAudioAmountPerArtist))...)
			}

		case false:
			resultItems = append(resultItems,
				c.collate(c.enquirer.getTopTracks(
					c.enquirer.getSimilarArtists(d.Artist, c.maxNumSimilarArtists),
					c.maxAudioAmountPerArtist))...)
		}
	}

	return
}

func (c collater) collateWithoutArtistStrain(data datastruct.AudioItems) []datastruct.AudioItem {
	if len(data.Items) > c.maxAudioAmountPerSource {
		return data.Items[:c.maxAudioAmountPerSource]
	}
	return data.Items
}

func (c collater) collateWithArtistStrain(data datastruct.AudioItems) []datastruct.AudioItem {
	sort.SliceStable(data.Items, func(i, j int) bool {
		return data.Items[i].Artist < data.Items[j].Artist
	})

	numberOfArtistSongs := make(map[string]int)
	var artistName string

	for i := 0; i < len(data.Items)-1; i++ {
		artistName = data.Items[i].Artist

		if artistName == data.Items[i+1].Artist {
			numberOfArtistSongs[artistName]++

			if numberOfArtistSongs[artistName] >= c.options.maxAudioAmountPerArtist {
				data.Items = append(data.Items[:i], data.Items[i+1:]...)
				i--
			}
		}
	}

	if len(data.Items) >= c.options.maxAudioAmountPerSource {
		return data.Items[:c.options.maxAudioAmountPerSource]
	}
	return data.Items
}
