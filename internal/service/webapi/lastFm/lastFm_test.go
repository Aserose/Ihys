package lastFm

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/repository"
	"IhysBestowal/pkg/customLogger"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLastFm(T *testing.T) {
	logs := customLogger.NewLogger()
	lfm := newTestLfm(logs)
	userId := int64(0)
	sourceItems := newSourceItems(
		newSong("Reliq", "gem"),
		newSong("Uniforms", "Serena"))

	convey.Convey("init", T, func() {

		convey.Convey("similar", func() { lfm.similar(userId, sourceItems) })
		convey.Convey("top", func() { lfm.top(sourceItems) })

	})
}

type testLfm struct {
	lfm ILastFM
}

func newTestLfm(log customLogger.Logger) testLfm {
	return testLfm{
		lfm: newLfm(log),
	}
}

func (t testLfm) similar(userId int64, sourceItems datastruct.AudioItems) {
	sim := func(amountPerSource int) {
		equalValue := amountPerSource * len(sourceItems.Items)
		assertion := convey.ShouldEqual
		if amountPerSource < 0 {
			equalValue = 0
		}
		if amountPerSource > 6 {
			equalValue = amountPerSource * 6
			assertion = convey.ShouldBeGreaterThanOrEqualTo
		}

		convey.So(len(t.lfm.GetSimiliarSongsFromLast(userId, sourceItems, SetMaxAudioAmountPerSource(amountPerSource)).Items),
			assertion, equalValue)
	}

	for _, num := range []int{6, 1, 0, -4, 2} {
		sim(num)
	}

}

func (t testLfm) top(sourceItems datastruct.AudioItems) {
	getTopTracks := func(numberOfTopSongs int) {
		equalValue := numberOfTopSongs * len(sourceItems.Items)
		assertion := convey.ShouldEqual
		if numberOfTopSongs < 0 {
			equalValue = 0
		}
		if numberOfTopSongs > 6 {
			equalValue = numberOfTopSongs * 6
			assertion = convey.ShouldBeGreaterThanOrEqualTo
		}

		convey.So(
			len(t.lfm.GetTopTracks(getAListOfArtists(sourceItems.Items), numberOfTopSongs).Items),
			assertion, equalValue,
		)
	}

	for _, num := range []int{1, 3, 5, -1} {
		getTopTracks(num)
	}
}

func newSong(artist, songTitle string) datastruct.AudioItem {
	return datastruct.AudioItem{
		Artist: artist,
		Title:  songTitle,
	}
}

func newSourceItems(songs ...datastruct.AudioItem) datastruct.AudioItems {
	return datastruct.AudioItems{
		Items: songs,
	}
}

func getAListOfArtists(items []datastruct.AudioItem) []string {
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = getArtist(item)
	}
	return result
}

func getArtist(item datastruct.AudioItem) string {
	return item.Artist
}

func newLfm(log customLogger.Logger) ILastFM {
	cfg := config.NewCfg(log)

	return NewLastFM(
		log,
		config.NewCfg(log).LastFM,
		repository.NewRepository(log, cfg.Repository))
}
