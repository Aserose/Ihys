package lastFm

import (
	"IhysBestowal/internal/datastruct"
	"sort"
	"sync"
)

type clt struct {
	opt
	enq
	collate func(data datastruct.Songs) []datastruct.Song
}

func newClt(enq enq, opts ...Set) clt {
	cl := clt{
		opt: opt{
			flowSize:          3,
			maxPerSource:      3,
			maxPerArtist:      1,
			maxSimilarArtists: 35,
			maxTopPerArtist:   10,
		},
		enq: enq,
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

func (c clt) SimilarParallel(uid int64, src datastruct.Songs) datastruct.Songs {
	if c.opt.maxPerSource <= 0 || c.opt.maxPerArtist <= 0 {
		return datastruct.Songs{}
	}

	res := []datastruct.Song{}
	wg := &sync.WaitGroup{}
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

	for low := 0; low <= len(src.Songs); low += c.opt.flowSize {
		if len(src.Songs[low:]) <= c.opt.flowSize {
			ch <- c.similar(src.Songs[low:])
			break
		}

		wg.Add(1)
		go func(source []datastruct.Song) {
			defer wg.Done()
			ch <- c.similar(source)
		}(src.Songs[low : low+c.opt.flowSize])
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
	for _, d := range src {
		sim := c.enq.similar(d.Artist, d.Title)

		switch sim.Songs != nil {
		case true:
			cltd := c.collate(sim)
			res = append(res, cltd...)

			if len(cltd) < c.maxPerSource {
				res = append(res, c.collate(
					c.enq.top(
						c.enq.similarArtists(d.Artist, c.maxPerSource-len(cltd)),
						c.maxPerArtist))...)
			}

		case false:
			res = append(res,
				c.collate(c.enq.top(
					c.enq.similarArtists(d.Artist, c.maxSimilarArtists),
					c.maxPerArtist))...)
		}
	}

	return
}

func (c clt) withoutArtistStrain(s datastruct.Songs) []datastruct.Song {
	if len(s.Songs) > c.maxPerSource {
		return s.Songs[:c.maxPerSource]
	}
	return s.Songs
}

func (c clt) withArtistStrain(s datastruct.Songs) []datastruct.Song {
	sort.SliceStable(s.Songs, func(i, j int) bool {
		return s.Songs[i].Artist < s.Songs[j].Artist
	})

	numArtistSongs := make(map[string]int)
	var artist string

	for i := 0; i < len(s.Songs)-1; i++ {
		artist = s.Songs[i].Artist

		if artist == s.Songs[i+1].Artist {
			numArtistSongs[artist]++

			if numArtistSongs[artist] >= c.opt.maxPerArtist {
				s.Songs = append(s.Songs[:i], s.Songs[i+1:]...)
				i--
			}
		}
	}

	if len(s.Songs) >= c.opt.maxPerSource {
		return s.Songs[:c.opt.maxPerSource]
	}
	return s.Songs
}
