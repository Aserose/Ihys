package webapi

import (
	"IhysBestowal/internal/service/webapi/lastFm"
	"IhysBestowal/internal/service/webapi/soundcloud"
	"IhysBestowal/internal/service/webapi/yaMusic"
)

const Frm = "all"

type From struct {
	soundcloud string
	yaM        string
	lastFm     LfmFrom
	allF       string
}

func newFrom() From {
	return From{
		soundcloud: soundcloud.From,
		yaM:        yaMusic.From,
		lastFm:     newLfmFrom(),
		allF:       Frm,
	}
}

func (s From) SoundCloud() string { return s.soundcloud }
func (s From) YaMusic() string    { return s.yaM }
func (s From) Lfm() LfmFrom       { return s.lastFm }
func (s From) All() string        { return s.allF }

type LfmFrom struct {
	Similar string
	Top     string
}

func newLfmFrom() LfmFrom {
	return LfmFrom{
		Similar: lastFm.From,
		Top:     lastFm.FromTop,
	}
}
