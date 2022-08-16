package soundcloud

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestParser(t *testing.T) {
	logs := customLogger.NewLogger()
	p := newTestParser(logs)

	defer p.CloseBrowser()

	convey.Convey("init", t, func() {

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

func (t testParser) sim() {
	convey.So(t.similar("Clark", "Winter Linn"), convey.ShouldNotEqual, []datastruct.Song{})
	convey.So(t.similar("gw324g2", "23g233r"), convey.ShouldResemble, []datastruct.Song{})
}
