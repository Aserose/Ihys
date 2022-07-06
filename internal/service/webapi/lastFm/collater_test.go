package lastFm

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
	"github.com/smartystreets/goconvey/convey"
	"sort"
	"testing"
)

func TestCollater(t *testing.T) {
	log := customLogger.NewLogger()
	cltr := newTestCollater(newTestEnquirer(log))
	userId := int64(0)
	sourceItems := newSourceItems(
		newSong("Reliq", "gem"),
		newSong("Uniforms", "Serena"),
		newSong("Telepopmusik", "Close"),
		newSong("Losers", "This Is A War"))

	convey.Convey("init", t, func() {

		convey.Convey("option <MaxAudioAmountPerSource>", func() { cltr.maxAudioAmountPerSource(userId, sourceItems) })
		convey.Convey("option <MaxAudioAmountPerArtist>", func() { cltr.maxAudioAmountPerArtist(userId, sourceItems) })

	})

}

type testCollater struct {
	enq enquirer
}

func newTestCollater(enq enquirer) testCollater {
	return testCollater{
		enq: enq,
	}
}

func (t testCollater) maxAudioAmountPerSource(userId int64, sourceItems datastruct.AudioItems) {
	getSimilar := func(maxAudioAmountPerSource int) {
		equalValue := maxAudioAmountPerSource * len(sourceItems.Items)
		assertion := convey.ShouldEqual
		if maxAudioAmountPerSource < 0 {
			equalValue = 0
		}
		if maxAudioAmountPerSource > 30 {
			equalValue = 30 * len(sourceItems.Items)
			assertion = convey.ShouldBeGreaterThanOrEqualTo
		}

		convey.So(
			len(newCollater(t.enq, SetMaxAudioAmountPerSource(maxAudioAmountPerSource)).getSimilarParallel(userId, sourceItems).Items),
			assertion,
			equalValue,
		)
	}

	for _, num := range []int{4, 74, 0, -4} {
		getSimilar(num)
	}
}

func (t testCollater) maxAudioAmountPerArtist(userId int64, sourceItems datastruct.AudioItems) {
	artistAlreadyOnTheList := func(listArtist []string) bool {
		counter := make(map[string]int)
		for _, artist := range listArtist {
			counter[artist]++
			if counter[artist] > 1 {
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

		listArtists := getAListOfArtists(newCollater(t.enq, SetMaxAudioAmountPerArtist(maxAudioAmountPerArtist)).getSimilarParallel(userId, sourceItems).Items)
		sort.Strings(listArtists)
		convey.So(artistAlreadyOnTheList(listArtists), convey.ShouldEqual, equalValue)
	}

	for _, num := range []int{1, 4} {
		testFunc(num)
	}
}
