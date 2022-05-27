package yaMusic

import (
	"IhysBestowal/pkg/customLogger"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCollater(t *testing.T) {
	log := customLogger.NewLogger()
	sourceItems := newSourceItems(newSong("boa", "duvet"))

	convey.Convey("init", t, func() {

		convey.Convey("song search", func() {
			convey.So(newCollater(log).getAudio("does 214 it offend you we are").Title, convey.ShouldEqual, "We Are Rockstars")

		})

		convey.Convey("similiar", func() {
			getSimiliar := func(maxAudioAmountPerSource int) {
				equalValue := maxAudioAmountPerSource*len(sourceItems.Items)
				if maxAudioAmountPerSource < 0 { equalValue = 0 }
				if maxAudioAmountPerSource > 10 { equalValue = 10}

				convey.So(
					len(newCollater(log, setMaxAudioAmountPerSource(maxAudioAmountPerSource)).getSimiliar(sourceItems).Items),
					convey.ShouldEqual, equalValue)
			}

			for _, num := range []int{ 5, 0, -3, 2532 }{
				getSimiliar(num)
			}
		})
	})
}
