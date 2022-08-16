package yaMusic

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestYaMusic(t *testing.T) {
	log := customLogger.NewLogger()
	ya := newTestYaMusic(log)
	src := newSrc(newSong("boa", "duvet"))

	convey.Convey(" ", t, func() {

		convey.Convey("find", func() { ya.find() })
		convey.Convey("similar", func() { ya.similar(src) })

	})

}

type testYaMusic struct {
	YaMusic
}

func newTestYaMusic(log customLogger.Logger) testYaMusic {
	return testYaMusic{
		YaMusic: New(log),
	}
}

func (t testYaMusic) find() {
	convey.So(t.Find("does 214 it offend you we are").Title, convey.ShouldEqual, "We Are Rockstars")
}

func (t testYaMusic) similar(sourceItems datastruct.Songs) {
	convey.So(len(t.Similar(sourceItems, MaxPerSource(3)).Songs), convey.ShouldEqual, 3)
}

func newSong(artist, songTitle string) datastruct.Song {
	return datastruct.Song{
		Artist: artist,
		Title:  songTitle,
	}
}

func newSrc(songs ...datastruct.Song) datastruct.Songs {
	return datastruct.Songs{
		Songs: songs,
	}
}
