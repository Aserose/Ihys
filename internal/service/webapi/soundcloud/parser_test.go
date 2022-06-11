package soundcloud

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestParser(t *testing.T) {
	logs := customLogger.NewLogger()
	parser := newTestParser(logs)

	defer parser.p.CloseBrowser()

	convey.Convey("init", t, func() {

		convey.Convey("similar", func() { parser.similar() })

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

func (t testParser) similar() {
	convey.So(t.p.getSimilar("Clark", "Winter Linn"), convey.ShouldNotEqual, []datastruct.AudioItem{})
	convey.So(t.p.getSimilar("gw324g2", "23g233r"), convey.ShouldResemble, []datastruct.AudioItem{})
}
