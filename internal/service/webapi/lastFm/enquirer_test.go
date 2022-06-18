package lastFm

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/repository"
	"IhysBestowal/pkg/customLogger"
	"github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestEnquirer(t *testing.T) {
	logs := customLogger.NewLogger()
	testEnq := enqStructWrap(newTestEnquirer(logs))

	artistList := []string{
		"Caspian", "Serph", "Akira Yamaoka",
	}
	query := "Suiri Taniuchi Hideki"

	convey.Convey("init", t, func() {

		convey.Convey("similar artists", func() { testEnq.testSimiliarArtists(artistList) })
		convey.Convey("top tracks", func() { testEnq.testTopTracks(artistList) })
		convey.Convey("song search", func() { testEnq.getAudio(query) })

	})
}

func newTestEnquirer(log customLogger.Logger) enquirer {
	cfg := config.NewCfg(log)

	return newEnquirer(log, cfg.LastFM, repository.NewRepository(log, cfg.Repository))
}

type testEnquirer struct {
	enq enquirer
}

func enqStructWrap(enq enquirer) testEnquirer {
	return testEnquirer{
		enq: enq,
	}
}

func (t testEnquirer) getAudio(query string) {
	convey.So(t.enq.getAudio(query), convey.ShouldNotResemble, datastruct.AudioItem{})
}

func (t testEnquirer) testTopTracks(artistList []string) {
	getTopTracks := func(num int) {
		equalValue := len(artistList) * num
		assertion := convey.ShouldEqual
		if num < 0 {
			equalValue = 0
		}
		if num > 4 {
			equalValue = len(artistList) * 4
			assertion = convey.ShouldBeGreaterThanOrEqualTo
		}

		convey.So(
			len(t.enq.getTopTracks(artistList, num).Items),
			assertion,
			equalValue,
		)
	}
	for _, num := range []int{1, 7, 0, -2} {
		getTopTracks(num)
	}
}

func (t testEnquirer) testSimiliarArtists(artistList []string) {
	getSimiliarArtists := func(limit int) {
		equalValue := len(artistList) * limit
		assertion := convey.ShouldEqual
		if limit < 0 {
			equalValue = 0
		}
		if limit > 4 {
			equalValue = len(artistList) * 4
			assertion = convey.ShouldBeGreaterThanOrEqualTo
		}

		for _, enumType := range enumTypes {
			convey.So(
				len(t.enq.getSimilarArtists(strings.Join(artistList, enumType), limit)),
				assertion,
				equalValue,
			)
		}
	}

	for _, limit := range []int{0, 16, 3, 8, -4} {
		getSimiliarArtists(limit)
	}
}
