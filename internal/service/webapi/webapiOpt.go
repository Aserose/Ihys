package webapi

import (
	"IhysBestowal/internal/service/webapi/lastFm"
	"IhysBestowal/internal/service/webapi/soundcloud"
	"IhysBestowal/internal/service/webapi/yaMusic"
)

type Opt struct {
	OneAudioPerArtist bool
	Ya                []yaMusic.ProcessingOptions
	Lf                []lastFm.ProcessingOptions
	Sc                []soundcloud.ProcessingOptions
}

func GetOptDefaultPreset() Opt {
	return Opt{
		OneAudioPerArtist: true,
		Ya: []yaMusic.ProcessingOptions{
			yaMusic.SetMaxAudioAmountPerSource(10),
		},
		Lf: []lastFm.ProcessingOptions{
			lastFm.SetMaxAudioAmountPerSource(200),
		},
		Sc: []soundcloud.ProcessingOptions{
			soundcloud.SetMaxAudioAmountPerSource(100),
		},
	}
}
