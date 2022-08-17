package yaMusic

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestParser(t *testing.T) {
	log := customLogger.NewLogger()
	p := newTestParser(log)

	convey.Convey(" ", t, func() {

		convey.Convey("find", func() { p.search() })
		convey.Convey("similar", func() { p.sim() })

	})
}

type testParser struct {
	parser
}

func newTestParser(log customLogger.Logger) testParser {
	return testParser{
		parser: newParser(log),
	}
}

func (t testParser) search() {
	convey.So(t.find("does 214 it offend you we are").Title, convey.ShouldEqual, "We Are Rockstars")
}

func (t testParser) sim() {
	convey.So(t.similar("Clark", "Winter Linn"), convey.ShouldNotEqual, []datastruct.Song{})
}
