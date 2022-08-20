package soundcloud

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSoundcloud(t *testing.T) {
	log := customLogger.NewLogger()
	sc := newTestSnCld(log)
	src := newSrc(newSong("boa", "duvet"))
	defer sc.Close()

	convey.Convey(" ", t, func() {

		convey.Convey("similar", func() { sc.similar(src) })

	})

}

type testSnCld struct {
	Soundcloud
}

func newTestSnCld(log customLogger.Logger) testSnCld {
	return testSnCld{
		Soundcloud: New(log),
	}
}

func (t testSnCld) similar(src datastruct.Set) {
	convey.So(len(t.Similar(src, MaxPerSource(3)).Song), convey.ShouldEqual, 3)
}

func newSong(artist, song string) datastruct.Song {
	return datastruct.Song{
		Artist: artist,
		Title:  song,
	}
}

func newSrc(s ...datastruct.Song) datastruct.Set {
	return datastruct.Set{
		Song: s,
	}
}
