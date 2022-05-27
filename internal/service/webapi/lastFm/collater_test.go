package lastFm

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/repository"
	"IhysBestowal/pkg/customLogger"
	"github.com/smartystreets/goconvey/convey"
	"sort"
	"testing"
)

func TestCollater(t *testing.T) {
	log := customLogger.NewLogger()
	enq := newTestEnquirer(log)
	sourceItems := newSourceItems(
		newSong("Reliq", "gem"),
		newSong("Uniforms", "Serena"))
	userId := int64(0)

	convey.Convey("init", t, func() {

		convey.Convey("option <MaxAudioAmountPerSource>", func() {
			getSimiliar := func(maxAudioAmountPerSource int) {
				equalValue := maxAudioAmountPerSource*len(sourceItems.Items)
				assertion := convey.ShouldEqual
				if maxAudioAmountPerSource < 0 { equalValue = 0 }
				if maxAudioAmountPerSource > 20 { assertion = convey.ShouldBeLessThanOrEqualTo }

				convey.So(
					len(newCollater(enq, setMaxAudioAmountPerSource(maxAudioAmountPerSource)).getSimiliars(userId, sourceItems).Items),
					assertion, equalValue)
			}

			for _, num := range []int{ 4, 74, 0 -4 } {
				getSimiliar(num)
			}

		})

		convey.Convey("option <MaxAudioAmountPerArtist>", func() {
			artistAlreadyOnTheList := func(listArtist []string) bool {
				counter := make(map[string]int)
				for _, artist := range listArtist {
					counter[artist]++
					if counter[artist] > 1 { return true }
				}
				return false
			}

			testFunc := func(maxAudioAmountPerSource int) {
				equalValue := false
				if maxAudioAmountPerSource > 1 { equalValue = true }

				listArtists := getAListOfArtists(newCollater(enq, setMaxAudioAmountPerArtist(maxAudioAmountPerSource)).getSimiliars(userId, sourceItems).Items)
				sort.Strings(listArtists)
				convey.So(artistAlreadyOnTheList(listArtists), convey.ShouldEqual, equalValue)
			}

			for _, num := range []int{ 1, 4 } {
				testFunc(num)
			}
		})
	})
}

func newTestEnquirer(log customLogger.Logger) iEnquirer {
	cfg := config.NewCfg(log)

	return newEnquirer(log, cfg.LastFM, repository.NewRepository(log, cfg.Repository))
}