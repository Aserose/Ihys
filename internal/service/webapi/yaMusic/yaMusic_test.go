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
	sourceItems := newSourceItems(newSong("boa","duvet"))

	convey.Convey("init", t, func() {

		convey.Convey("song search", func() { ya.songSearch() })
		convey.Convey("similiar", func() { ya.similiar(sourceItems) })

	})

}

type testYaMusic struct {
	ya IYaMusic
}

func newTestYaMusic(log customLogger.Logger) testYaMusic{
	return testYaMusic{
		ya: NewYaMusic(log),
	}
}

func (t testYaMusic) songSearch() {
	convey.So(t.ya.GetAudio("does 214 it offend you we are").Title, convey.ShouldEqual, "We Are Rockstars")
}

func (t testYaMusic) similiar(sourceItems datastruct.AudioItems) {
	convey.So(len(t.ya.GetSimliarSongsFromYa10(sourceItems).Items), convey.ShouldEqual, 10)
}

func newSong(artist, songTitle string) datastruct.AudioItem{
	return datastruct.AudioItem{
		Artist: artist,
		Title: songTitle,
	}
}

func newSourceItems(songs ... datastruct.AudioItem) datastruct.AudioItems{
	return datastruct.AudioItems{
		Items: songs,
	}
}