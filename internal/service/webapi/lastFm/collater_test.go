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
		newSong("Uniforms", "Serena"))

	convey.Convey("init", t, func() {

		convey.Convey("option <MaxAudioAmountPerSource>", func() { cltr.maxAudioAmountPerSource(userId, sourceItems) })
		convey.Convey("option <MaxAudioAmountPerArtist>", func() { cltr.maxAudioAmountPerArtist(userId, sourceItems) })

	})

}

type testCollater struct {
	enq iEnquirer
}

func newTestCollater(enq iEnquirer) testCollater {
	return testCollater{
		enq: enq,
	}
}

func (t testCollater) maxAudioAmountPerSource(userId int64, sourceItems datastruct.AudioItems) {
	getSimiliar := func(maxAudioAmountPerSource int) {
		equalValue := []interface{}{maxAudioAmountPerSource * len(sourceItems.Items)}
		assertion := convey.ShouldEqual
		if maxAudioAmountPerSource < 0 {
			equalValue = []interface{}{0}
		}
		if maxAudioAmountPerSource > 20 {
			equalValue = []interface{}{20 * len(sourceItems.Items)}
			assertion = convey.ShouldBeGreaterThanOrEqualTo
		}

		convey.So(
			len(newCollater(t.enq, setMaxAudioAmountPerSource(maxAudioAmountPerSource)).getSimiliars(userId, sourceItems).Items),
			assertion,
			equalValue...,
		)
	}

	for _, num := range []int{4, 74, 0, -4} {
		getSimiliar(num)
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

	testFunc := func(maxAudioAmountPerSource int) {
		equalValue := false
		if maxAudioAmountPerSource > 1 {
			equalValue = true
		}

		listArtists := getAListOfArtists(newCollater(t.enq, setMaxAudioAmountPerArtist(maxAudioAmountPerSource)).getSimiliars(userId, sourceItems).Items)
		sort.Strings(listArtists)
		convey.So(artistAlreadyOnTheList(listArtists), convey.ShouldEqual, equalValue)
	}

	for _, num := range []int{1, 4} {
		testFunc(num)
	}
}
