package soundcloud

import (
	"IhysBestowal/internal/datastruct"
	"sync"
)

type clt struct {
	parser
	opt
	collate func(data []datastruct.Song) []datastruct.Song
}

func newClt(p parser, opts ...Set) clt {
	cltr := clt{
		parser: p,
		opt: opt{
			FlowSize:     1,
			maxPerSource: 3,
			maxPerArtist: 1,
		},
	}

	if opts != nil {
		for _, set := range opts {
			set(&cltr.opt)
		}
	}

	if cltr.opt.maxPerArtist == 0 {
		cltr.collate = cltr.withoutArtistStrain
	} else {
		cltr.collate = cltr.withArtistStrain
	}

	return cltr
}

func (c clt) similarParallel(src datastruct.Songs) datastruct.Songs {
	if c.opt.maxPerSource <= 0 {
		return datastruct.Songs{}
	}

	wg := &sync.WaitGroup{}
	res := []datastruct.Song{}
	ch := make(chan []datastruct.Song)
	cls := make(chan struct{})

	go func() {
		for {
			select {
			case sim, ok := <-ch:
				if !ok {
					continue
				}
				res = append(res, sim...)
			case <-cls:
				return
			}
		}
	}()

	for low := 0; low <= len(src.Songs); low += c.opt.FlowSize {
		if len(src.Songs[low:]) <= c.opt.FlowSize {
			ch <- c.similar(src.Songs[low:])
			break
		}

		wg.Add(1)
		go func(s []datastruct.Song) {
			defer wg.Done()
			ch <- c.similar(s)
		}(src.Songs[low : low+c.opt.FlowSize])
	}

	wg.Wait()
	close(ch)
	cls <- struct{}{}
	close(cls)

	return datastruct.Songs{
		Songs: res,
		From:  From,
	}
}

func (c clt) similar(src []datastruct.Song) (res []datastruct.Song) {
	for _, s := range src {
		res = append(res, c.collate(c.parser.similar(s.Artist, s.Title))...)
	}
	return
}

func (c clt) withoutArtistStrain(data []datastruct.Song) []datastruct.Song {
	if len(data) >= c.opt.maxPerSource {
		return data[:c.opt.maxPerSource]
	}
	return data
}

func (c clt) withArtistStrain(data []datastruct.Song) []datastruct.Song {
	numArtistSongs := make(map[string]int)
	var artistName string

	for i := 0; i < len(data)-1; i++ {
		artistName = data[i].Artist

		if artistName == data[i+1].Artist {
			numArtistSongs[artistName]++

			if numArtistSongs[artistName] >= c.opt.maxPerArtist {
				data = append(data[:i], data[i+1:]...)
				i--
			}
		}
	}

	if len(data) >= c.opt.maxPerSource {
		return data[:c.opt.maxPerSource]
	}
	return data
}