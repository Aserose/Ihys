package lastFm

import (
	"IhysBestowal/internal/datastruct"
	"sort"
	"sync"
)

type clt struct {
	opt
	enq
	collate func(data datastruct.Set) []datastruct.Song
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

	for _, set := range opts {
		set(&cl.opt)
	}

	if cl.opt.maxPerArtist == 0 {
		cl.collate = cl.withoutArtistStrain
	} else {
		cl.collate = cl.withArtistStrain
	}

	return cl
}

func (c clt) SimilarParallel(uid int64, src datastruct.Set) datastruct.Set {
	if c.opt.maxPerSource <= 0 || c.opt.maxPerArtist <= 0 {
		return datastruct.Set{}
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

	for low := 0; low <= len(src.Song); low += c.opt.flowSize {
		if len(src.Song[low:]) <= c.opt.flowSize {
			ch <- c.similar(src.Song[low:])
			break
		}

		wg.Add(1)
		go func(source []datastruct.Song) {
			defer wg.Done()
			ch <- c.similar(source)
		}(src.Song[low : low+c.opt.flowSize])
	}

	wg.Wait()
	cls <- struct{}{}
	close(cls)
	close(ch)

	return datastruct.Set{
		Song: res,
		From: From,
	}
}

func (c clt) similar(src []datastruct.Song) (res []datastruct.Song) {
	for _, d := range src {
		sim := c.enq.similar(d.Artist, d.Title)

		switch sim.Song != nil {
		case true:
			cltd := c.collate(sim)
			res = append(res, cltd...)

			if len(cltd) < c.maxPerSource {
				res = append(res, c.collate(c.enq.top(c.enq.similarArtists(d.Artist, c.maxPerSource-len(cltd)), c.maxPerArtist))...)
			}

		case false:
			res = append(res, c.collate(c.enq.top(c.enq.similarArtists(d.Artist, c.maxSimilarArtists), c.maxPerArtist))...)
		}
	}

	return
}

func (c clt) withoutArtistStrain(s datastruct.Set) []datastruct.Song {
	if len(s.Song) > c.maxPerSource {
		return s.Song[:c.maxPerSource]
	}
	return s.Song
}

func (c clt) withArtistStrain(s datastruct.Set) []datastruct.Song {
	sort.SliceStable(s.Song, func(i, j int) bool {
		return s.Song[i].Artist < s.Song[j].Artist
	})

	var (
		numArtistSongs = make(map[string]int)
		artist         string
	)

	for i := 0; i < len(s.Song)-1; i++ {
		artist = s.Song[i].Artist

		if artist == s.Song[i+1].Artist {
			numArtistSongs[artist]++

			if numArtistSongs[artist] >= c.opt.maxPerArtist {
				s.Song = append(s.Song[:i], s.Song[i+1:]...)
				i--
			}
		}
	}

	if len(s.Song) >= c.opt.maxPerSource {
		return s.Song[:c.opt.maxPerSource]
	}
	return s.Song
}
