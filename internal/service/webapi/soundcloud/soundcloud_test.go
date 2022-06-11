package soundcloud

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSoundcloud(t *testing.T) {
	log := customLogger.NewLogger()
	sc := newTestSoundcloud(log)
	sourceItems := newSourceItems(newSong("boa", "duvet"))
	defer sc.sc.CloseBrowser()

	convey.Convey("init", t, func() {

		convey.Convey("similar", func() { sc.similar(sourceItems) })

	})

}

type testSoundcloud struct {
	sc ISoundcloud
}

func newTestSoundcloud(log customLogger.Logger) testSoundcloud {
	return testSoundcloud{
		sc: NewSoundcloud(log),
	}
}

func (t testSoundcloud) similar(sourceItems datastruct.AudioItems) {
	convey.So(len(t.sc.GetSimilar(sourceItems, SetMaxAudioAmountPerSource(3)).Items), convey.ShouldEqual, 3)
}

func newSong(artist, songTitle string) datastruct.AudioItem {
	return datastruct.AudioItem{
		Artist: artist,
		Title:  songTitle,
	}
}

func newSourceItems(songs ...datastruct.AudioItem) datastruct.AudioItems {
	return datastruct.AudioItems{
		Items: songs,
	}
}
