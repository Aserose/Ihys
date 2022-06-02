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
	sourceItems := newSourceItems(
		newSong("boa", "duvet"),
		newSong("Rogue Valley", "The Wolves & the Ravens"),
		newSong("Telepopmusik", "Close"))

	convey.Convey("init", t, func() {

		convey.Convey("song search", func() { c.songSearch() })
		convey.Convey("similiar", func() { c.similiar(sourceItems) })

	})
}

type testCollater struct {
	newCollater func(opts ...ProcessingOptions) collater
}

func newTestCollater(log customLogger.Logger) testCollater {
	parser := newParser(log)
	return testCollater{
		newCollater: func(opts ...ProcessingOptions) collater {
			return newCollater(parser, opts...)
		},
	}
}

func (t testCollater) songSearch() {
	convey.So(t.newCollater().getAudio("does 214 it offend you we are").Title, convey.ShouldEqual, "We Are Rockstars")
}

func (t testCollater) similiar(sourceItems datastruct.AudioItems) {
	getSimiliar := func(maxAudioAmountPerSource int) {
		equalValue := maxAudioAmountPerSource * len(sourceItems.Items)
		assertion := convey.ShouldEqual
		if maxAudioAmountPerSource < 0 {
			equalValue = 0
		}
		if maxAudioAmountPerSource > 20 {
			equalValue = 20
			assertion = convey.ShouldBeGreaterThanOrEqualTo
		}

		convey.So(
			len(t.newCollater(SetMaxAudioAmountPerSource(maxAudioAmountPerSource)).getSimilarParallel(sourceItems).Items),
			assertion, equalValue)
	}

	for _, num := range []int{5, 0, -3, 2532} {
		getSimiliar(num)
	}
}
