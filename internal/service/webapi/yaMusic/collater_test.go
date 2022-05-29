package yaMusic

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCollater(t *testing.T) {
	log := customLogger.NewLogger()
	c := newTestCollater(log)
	sourceItems := newSourceItems(newSong("boa", "duvet"))

	convey.Convey("init", t, func() {

		convey.Convey("song search", func() { c.songSearch() })
		convey.Convey("similiar", func() { c.similiar(sourceItems) })

	})
}

type testCollater struct {
	newCollater func(opts ...processingOptions) iCollater
}

func newTestCollater(log customLogger.Logger) testCollater {
	return testCollater{
		newCollater: func(opts ...processingOptions) iCollater {
			return newCollater(log, opts...)
		},
	}
}

func (t testCollater) songSearch() {
	convey.So(t.newCollater().getAudio("does 214 it offend you we are").Title, convey.ShouldEqual, "We Are Rockstars")
}

func (t testCollater) similiar(sourceItems datastruct.AudioItems) {
	getSimiliar := func(maxAudioAmountPerSource int) {
		equalValue := maxAudioAmountPerSource*len(sourceItems.Items)
		if maxAudioAmountPerSource < 0 { equalValue = 0 }
		if maxAudioAmountPerSource > 10 { equalValue = 10}

		convey.So(
			len(t.newCollater(setMaxAudioAmountPerSource(maxAudioAmountPerSource)).getSimiliar(sourceItems).Items),
			convey.ShouldEqual, equalValue)
	}

	for _, num := range []int{ 5, 0, -3, 2532 }{
		getSimiliar(num)
	}
}