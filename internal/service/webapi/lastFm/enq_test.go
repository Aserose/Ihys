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
	te := newTestEnq(logs)

	arts := []string{
		"Caspian", "Serph", "Akira Yamaoka",
	}
	query := "Suiri Taniuchi Hideki"

	convey.Convey(" ", t, func() {

		convey.Convey("similar artists", func() { te.similarArtists(arts) })
		convey.Convey("top tracks", func() { te.top(arts) })
		convey.Convey("song search", func() { te.search(query) })

	})
}

type testEnq struct {
	enq enq
}

func newTestEnq(log customLogger.Logger) testEnq {
	cfg := config.New(log)
	return testEnq{
		enq: newEnq(log, cfg.LastFM, repository.New(log, cfg.Repository)),
	}
}

func (t testEnq) search(query string) {
	convey.So(t.enq.find(query), convey.ShouldNotResemble, datastruct.Song{})
}

func (t testEnq) top(artists []string) {
	top := func(num int) {
		equalValue := len(artists) * num
		assertion := convey.ShouldEqual
		if num < 0 {
			equalValue = 0
		}
		if num > 4 {
			equalValue = len(artists) * 4
			assertion = convey.ShouldBeGreaterThanOrEqualTo
		}

		convey.So(len(t.enq.top(artists, num).Song), assertion, equalValue)
	}

	for _, num := range []int{1, 7, 0, -2} {
		top(num)
	}

}

func (t testEnq) similarArtists(artists []string) {
	similar := func(max int) {
		equalValue := len(artists) * max
		assertion := convey.ShouldEqual
		if max < 0 {
			equalValue = 0
		}
		if max > 4 {
			equalValue = len(artists) * 4
			assertion = convey.ShouldBeGreaterThanOrEqualTo
		}

		for _, enumType := range enumTypes {
			convey.So(
				len(t.enq.similarArtists(strings.Join(artists, enumType), max)), assertion, equalValue,
			)
		}
	}

	for _, max := range []int{0, 16, 3, 8, -4} {
		similar(max)
	}
}
