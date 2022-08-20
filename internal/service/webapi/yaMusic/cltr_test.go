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

		convey.Convey("similar", func() { c.similar(src) })
		convey.Convey("find", func() { c.find() })

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

func (t testClt) similar(src datastruct.Set) {
	get := func(maxPerSource int) {
		equalValue := maxPerSource * len(src.Song)
		assertion := convey.ShouldEqual
		if maxPerSource < 0 {
			equalValue = 0
		}
		if maxPerSource > 20 {
			equalValue = 20
			assertion = convey.ShouldBeGreaterThanOrEqualTo
		}

		convey.So(
			len(t.newClt(MaxPerSource(maxPerSource)).similarParallel(src).Song),
			assertion, equalValue)
	}

	for _, num := range []int{5, 0, -3, 2532} {
		get(num)
	}
}
