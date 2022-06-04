package webapi

import (
	"IhysBestowal/internal/service/webapi/lastFm"
	"IhysBestowal/internal/service/webapi/yaMusic"
)

type Opt struct {
	oneAudioPerArtist bool
	ya                []yaMusic.ProcessingOptions
	lf                []lastFm.ProcessingOptions
}

func GetOptDefaultPreset() Opt {
	return Opt{
		ya: []yaMusic.ProcessingOptions{
			yaMusic.SetMaxAudioAmountPerSource(10),
		},
		lf: []lastFm.ProcessingOptions{
			lastFm.SetMaxAudioAmountPerSource(200),
		},
	}
}
