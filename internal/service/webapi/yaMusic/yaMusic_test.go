package yaMusic

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestYaMusic(t *testing.T) {
	log := customLogger.NewLogger()
	ya := NewYaMusic(log)
	sourceItems := newSourceItems(newSong("boa","duvet"))

	convey.Convey("init", t, func() {

		convey.Convey("song search", func() {
			convey.So(ya.GetAudio("does 214 it offend you we are").Title, convey.ShouldEqual, "We Are Rockstars")

		})

		convey.Convey("similiar", func() {
			convey.So(len(ya.GetSimliarSongsFromYa10(sourceItems).Items), convey.ShouldEqual, 10)

		})

	})

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