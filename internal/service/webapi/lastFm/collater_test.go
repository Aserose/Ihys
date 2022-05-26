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

		convey.Convey("MaxAudioAmountPerSource", func() {
			maxAudioAmountPerSource := 4
			convey.So(
				len(newCollater(enq, setMaxAudioAmountPerSource(maxAudioAmountPerSource)).getSimiliars(userId, sourceItems).Items),
				convey.ShouldBeLessThanOrEqualTo,
				maxAudioAmountPerSource*len(sourceItems.Items))
		})

		convey.Convey("MaxAudioAmountPerArtist", func() {
			listArtists := getAListOfArtists(newCollater(enq, setMaxAudioAmountPerArtist(4)).getSimiliars(userId, sourceItems).Items)
			sort.Strings(listArtists)
			convey.So(listArtists[0], convey.ShouldEqual, listArtists[1])
		})
	})
}

func newTestEnquirer(log customLogger.Logger) iEnquirer {
	cfg := config.NewCfg(log)

	return newEnquirer(log, cfg.LastFM, repository.NewRepository(log, cfg.Repository))
}