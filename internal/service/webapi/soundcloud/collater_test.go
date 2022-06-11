package soundcloud

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCollater(t *testing.T) {
	logs := customLogger.NewLogger()
	c := newTestCollater(logs)
	sourceItems := newSourceItems(
		newSong("Neutral", "Next to the Stars"),
		newSong("Oh Hiroshima", "Mirage"),
		newSong("Серебряная Свадьба", "Ag"))

	defer c.closeBrowser()

	convey.Convey("init", t, func() {

		convey.Convey("similar", func() { c.similar(sourceItems) })

	})
}

type testCollater struct {
	newCollater  func(opts ...ProcessingOptions) collater
	closeBrowser func()
}

func newTestCollater(log customLogger.Logger) testCollater {
	parser := newParser(log)
	return testCollater{
		newCollater: func(opts ...ProcessingOptions) collater {
			return newCollater(parser, opts...)
		},
		closeBrowser: parser.CloseBrowser,
	}
}

func (t testCollater) shutdown() {
	t.closeBrowser()
}

func (t testCollater) similar(sourceItems datastruct.AudioItems) {
	getSimilar := func(maxAudioAmountPerSource int) {
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
		getSimilar(num)
	}
}
