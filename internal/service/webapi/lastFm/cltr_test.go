package lastFm

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/repository"
	"IhysBestowal/pkg/customLogger"
	"github.com/smartystreets/goconvey/convey"
	"sort"
	"testing"
)

func TestCollater(t *testing.T) {
	cl := newTestClt(customLogger.NewLogger())
	uid := int64(0)
	src := newSourceItems(
		newSong("Reliq", "gem"),
		newSong("Uniforms", "Serena"),
		newSong("Telepopmusik", "Close"),
		newSong("Losers", "This Is A War"))

	convey.Convey(" ", t, func() {

		convey.Convey("option <MaxAudioAmountPerSource>", func() { cl.maxPerSource(uid, src) })
		convey.Convey("option <MaxAudioAmountPerArtist>", func() { cl.maxPerArtist(uid, src) })

	})

}

type testClt struct {
	enq
}

func newTestClt(log customLogger.Logger) testClt {
	cfg := config.New(log)
	return testClt{
		enq: newEnq(log, cfg.LastFM, repository.New(log, cfg.Repository)),
	}
}

func (t testClt) maxPerSource(uid int64, src datastruct.Songs) {
	get := func(maxPerSource int) {
		equalValue := maxPerSource * len(src.Songs)
		assertion := convey.ShouldEqual
		if maxPerSource < 0 {
			equalValue = 0
		}
		if maxPerSource > 30 {
			equalValue = 30 * len(src.Songs)
			assertion = convey.ShouldBeGreaterThanOrEqualTo
		}

		convey.So(
			len(newClt(t.enq, MaxPerSource(maxPerSource)).SimilarParallel(uid, src).Songs),
			assertion,
			equalValue,
		)
	}

	for _, num := range []int{4, 74, 0, -4} {
		get(num)
	}
}

func (t testClt) maxPerArtist(userId int64, src datastruct.Songs) {
	isUniq := func(artists []string) bool {
		counter := make(map[string]int)
		for _, a := range artists {
			counter[a]++
			if counter[a] > 1 {
				return true
			}
		}
		return false
	}

	testFunc := func(maxAudioAmountPerArtist int) {
		equalValue := false
		if maxAudioAmountPerArtist > 1 {
			equalValue = true
		}

		arts := artists(newClt(t.enq, MaxPerArtist(maxAudioAmountPerArtist)).SimilarParallel(userId, src).Songs)
		sort.Strings(arts)
		convey.So(isUniq(arts), convey.ShouldEqual, equalValue)
	}

	for _, num := range []int{1, 4} {
		testFunc(num)
	}
}