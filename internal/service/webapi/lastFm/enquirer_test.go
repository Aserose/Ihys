package lastFm

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/repository"
	"IhysBestowal/pkg/customLogger"
	"github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestEnquirer(t *testing.T) {
	log := customLogger.NewLogger()
	testEnq := enqStructWrap(newTestEnquirer(log))
	artistList := []string{
		"Caspian", "Serph", "Akira Yamaoka",
	}

	convey.Convey("init", t, func() {

		convey.Convey("similar artists", func() { testEnq.testSimiliarArtists(artistList) })
		convey.Convey("top tracks", func() { testEnq.testTopTracks(artistList) })

	})
}

func newTestEnquirer(log customLogger.Logger) iEnquirer {
	cfg := config.NewCfg(log)

	return newEnquirer(log, cfg.LastFM, repository.NewRepository(log, cfg.Repository))
}

type testEnquirer struct {
	enq iEnquirer
}

func enqStructWrap(enq iEnquirer) testEnquirer {
	return testEnquirer{
		enq: enq,
	}
}

func (t testEnquirer) testTopTracks(artistList []string) {
	getTopTracks := func(num int) {
		equalValue := []interface{}{len(artistList) * num}
		assertion := convey.ShouldEqual
		if num < 0 {
			equalValue = []interface{}{0}
		}
		if num > 4 {
			equalValue = []interface{}{len(artistList) * 4}
			assertion = convey.ShouldBeGreaterThanOrEqualTo
		}

		convey.So(
			len(t.enq.getTopTracks(artistList, num).Items),
			assertion,
			equalValue...,
		)
	}
	for _, num := range []int{1, 7, 0, -2} {
		getTopTracks(num)
	}
}

func (t testEnquirer) testSimiliarArtists(artistList []string) {
	getSimiliarArtists := func(limit int) {
		equalValue := []interface{}{len(artistList) * limit}
		assertion := convey.ShouldEqual
		if limit < 0 {
			equalValue = []interface{}{0}
		}
		if limit > 4 {
			equalValue = []interface{}{len(artistList) * 4}
			assertion = convey.ShouldBeGreaterThanOrEqualTo
		}

		for _, enumType := range enumTypes {
			convey.So(
				len(t.enq.getSimilarArtists(strings.Join(artistList, enumType), limit)),
				assertion,
				equalValue...,
			)
		}
	}

	for _, limit := range []int{0, 3, 8, -4} {
		getSimiliarArtists(limit)
	}
}