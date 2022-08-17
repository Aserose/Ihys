package yaMusic

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCollater(t *testing.T) {
	log := customLogger.NewLogger()
	c := newTestClt(log)
	src := newSrc(
		newSong("boa", "duvet"),
		newSong("Rogue Valley", "The Wolves & the Ravens"),
		newSong("Telepopmusik", "Close"))

	convey.Convey(" ", t, func() {

		convey.Convey("find", func() { c.find() })
		convey.Convey("similar", func() { c.similar(src) })

	})
}

type testClt struct {
	newClt func(opts ...Set) clt
}

func newTestClt(log customLogger.Logger) testClt {
	p := newParser(log)
	return testClt{
		newClt: func(opts ...Set) clt {
			return newClt(p, opts...)
		},
	}
}

func (t testClt) find() {
	convey.So(t.newClt().find("does 214 it offend you we are").Title, convey.ShouldEqual, "We Are Rockstars")
}

func (t testClt) similar(sourceItems datastruct.Songs) {
	get := func(maxAudioAmountPerSource int) {
		equalValue := maxAudioAmountPerSource * len(sourceItems.Songs)
		assertion := convey.ShouldEqual
		if maxAudioAmountPerSource < 0 {
			equalValue = 0
		}
		if maxAudioAmountPerSource > 20 {
			equalValue = 20
			assertion = convey.ShouldBeGreaterThanOrEqualTo
		}

		convey.So(
			len(t.newClt(MaxPerSource(maxAudioAmountPerSource)).similarParallel(sourceItems).Songs),
			assertion, equalValue)
	}

	for _, num := range []int{5, 0, -3, 2532} {
		get(num)
	}
}
