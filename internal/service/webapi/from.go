package webapi

import (
	"IhysBestowal/internal/service/webapi/lastFm"
	"IhysBestowal/internal/service/webapi/soundcloud"
	"IhysBestowal/internal/service/webapi/yaMusic"
)

const Frm = "all"

type From struct {
	scd  string
	ya   string
	lfm  LfmFrom
	allF string
}

func newFrom() From {
	return From{
		scd:  soundcloud.From,
		ya:   yaMusic.From,
		lfm:  newLfmFrom(),
		allF: Frm,
	}
}

func (s From) SoundCloud() string { return s.scd }
func (s From) YaMusic() string    { return s.ya }
func (s From) LastFm() LfmFrom    { return s.lfm }
func (s From) All() string        { return s.allF }

type LfmFrom struct {
	similar string
	top     string
}

func newLfmFrom() LfmFrom {
	return LfmFrom{
		similar: lastFm.From,
		top:     lastFm.FromTop,
	}
}

func (l LfmFrom) Similar() string { return l.similar }
func (l LfmFrom) Top() string     { return l.top }
