package yaMusic

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)


func TestParser(t *testing.T) {
	log := customLogger.NewLogger()
	parser := newParser(log)

	convey.Convey("init", t, func() {

		convey.Convey("song search", func() {
			convey.So(parser.getAudio("does 214 it offend you we are").Title, convey.ShouldEqual, "We Are Rockstars")

		})

		convey.Convey("similiar", func() {
			convey.So(parser.getSimiliars("Clark", "Winter Linn"), convey.ShouldNotEqual, datastruct.YaMSimiliar{})
		})
	})
}