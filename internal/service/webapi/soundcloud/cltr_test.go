package soundcloud

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestClt(t *testing.T) {
	logs := customLogger.NewLogger()
	c := newTestClt(logs)
	src := newSrc(
		newSong("Neutral", "Next to the Stars"),
		newSong("Oh Hiroshima", "Mirage"),
		newSong("Серебряная Свадьба", "Ag"))

	defer c.closeBrowser()

	convey.Convey(" ", t, func() {

		convey.Convey("similar", func() { c.similar(src) })

	})
}

type testClt struct {
	newClt       func(s ...Set) clt
	closeBrowser func()
}

func newTestClt(log customLogger.Logger) testClt {
	p := newParser(log)
	return testClt{
		newClt: func(opts ...Set) clt {
			return newClt(p, opts...)
		},
		closeBrowser: p.CloseBrowser,
	}
}

func (t testClt) shutdown() {
	t.closeBrowser()
}

func (t testClt) similar(src datastruct.Songs) {
	get := func(maxPerSource int) {
		equalValue := maxPerSource * len(src.Songs)
		assertion := convey.ShouldEqual
		if maxPerSource < 0 {
			equalValue = 0
		}
		if maxPerSource > 10 {
			equalValue = 10
			assertion = convey.ShouldBeGreaterThanOrEqualTo
		}

		convey.So(
			len(t.newClt(MaxPerSource(maxPerSource)).similarParallel(src).Songs), assertion, equalValue)
	}

	for _, num := range []int{5, 0, -3, 2532} {
		get(num)
	}
}
