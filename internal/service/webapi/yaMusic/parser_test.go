package yaMusic

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestParser(t *testing.T) {
	log := customLogger.NewLogger()
	parser := newTestParser(log)

	convey.Convey("init", t, func() {

		convey.Convey("song search", func() { parser.songSearch() })
		convey.Convey("similiar", func() { parser.similar() })

	})
}

type testParser struct {
	p parser
}

func newTestParser(log customLogger.Logger) testParser {
	return testParser{
		p: newParser(log),
	}
}

func (t testParser) songSearch() {
	convey.So(t.p.getAudio("does 214 it offend you we are").Title, convey.ShouldEqual, "We Are Rockstars")
}

func (t testParser) similar() {
	convey.So(t.p.getSimilar("Clark", "Winter Linn"), convey.ShouldNotEqual, []datastruct.AudioItem{})
}
