package webapi

import (
	"IhysBestowal/internal/service/webapi/lastFm"
	"IhysBestowal/internal/service/webapi/soundcloud"
	"IhysBestowal/internal/service/webapi/yaMusic"
)

type Opt struct {
	OnePerArtist bool
	Ya           []yaMusic.Set
	Lf           []lastFm.Set
	Sc           []soundcloud.Set
}

func Default() Opt {
	return Opt{
		OnePerArtist: true,
		Ya: []yaMusic.Set{
			yaMusic.MaxPerSource(10),
		},
		Lf: []lastFm.Set{
			lastFm.MaxPerSource(200),
		},
		Sc: []soundcloud.Set{
			soundcloud.MaxPerSource(100),
		},
	}
}
