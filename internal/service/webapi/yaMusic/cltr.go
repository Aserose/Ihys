package yaMusic

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
	cl := clt{
		parser: p,
		opt: opt{
			flowSize:     3,
			maxPerSource: 3,
			maxPerArtist: 1,
		},
	}

	if opts != nil {
		for _, set := range opts {
			set(&cl.opt)
		}
	}

	if cl.opt.maxPerArtist == 0 {
		cl.collate = cl.withoutArtistStrain
	} else {
		cl.collate = cl.withArtistStrain
	}

	return cl
}

func (m clt) similarParallel(src datastruct.Songs) datastruct.Songs {
	if m.opt.maxPerSource <= 0 {
		return datastruct.Songs{}
	}

	wg := &sync.WaitGroup{}
	res := []datastruct.Song{}
	ch := make(chan []datastruct.Song)
	cls := make(chan struct{})

	go func() {
		for {
			select {
			case sim := <-ch:
				res = append(res, sim...)
			case <-cls:
				return
			}
		}
	}()

	for low := 0; low <= len(src.Songs); low += m.opt.flowSize {
		if len(src.Songs[low:]) <= m.opt.flowSize {
			ch <- m.similar(src.Songs[low:])
			break
		}

		wg.Add(1)
		go func(source []datastruct.Song) {
			defer wg.Done()
			ch <- m.similar(source)
		}(src.Songs[low : low+m.opt.flowSize])

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

func (m clt) similar(src []datastruct.Song) (res []datastruct.Song) {
	for _, item := range src {
		res = append(res, m.collate(m.parser.similar(item.Artist, item.Title))...)
	}
	return
}

func (m clt) withoutArtistStrain(s []datastruct.Song) []datastruct.Song {
	if len(s) >= m.opt.maxPerSource {
		return s[:m.opt.maxPerSource]
	}
	return s
}

func (m clt) withArtistStrain(s []datastruct.Song) []datastruct.Song {
	numArtistSongs := make(map[string]int)
	var artist string

	for i := 0; i < len(s)-1; i++ {
		artist = s[i].Artist

		if artist == s[i+1].Artist {
			numArtistSongs[artist]++

			if numArtistSongs[artist] >= m.opt.maxPerArtist {
				s = append(s[:i], s[i+1:]...)
				i--
			}
		}
	}

	if len(s) >= m.opt.maxPerSource {
		return s[:m.opt.maxPerSource]
	}
	return s
}
