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
	log := customLogger.NewLogger()
	lfm := newTestLfm(log)
	userId := int64(0)
	sourceItems := newSourceItems(
		newSong("Reliq", "gem"),
		newSong("Uniforms", "Serena"))

	convey.Convey("init", T, func() {

		convey.Convey("similar", func() { lfm.similiar(userId, sourceItems) })
		convey.Convey("similiar100", func() { lfm.similiar100(userId, sourceItems) })
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

func (t testLfm) similiar(userId int64, sourceItems datastruct.AudioItems) {
	convey.So(
		t.lfm.GetSimiliarSongsFromLast(userId, sourceItems),
		convey.ShouldNotEqual,
		datastruct.AudioItems{})
}

func (t testLfm) similiar100(userId int64, sourceItems datastruct.AudioItems) {
	convey.So(
		len(t.lfm.GetSimiliarSongsFromLast100(userId, sourceItems).Items),
		convey.ShouldBeGreaterThan,
		3*len(sourceItems.Items))
}

func (t testLfm) top(sourceItems datastruct.AudioItems) {
	getTopTracks := func(numberOfTopSongs int) {
		equalValue := []interface{}{numberOfTopSongs * len(sourceItems.Items)}
		assertion := convey.ShouldEqual
		if numberOfTopSongs < 0 {
			equalValue = []interface{}{0}
		}
		if numberOfTopSongs > 6 {
			equalValue = []interface{}{numberOfTopSongs * 6}
			assertion = convey.ShouldBeGreaterThanOrEqualTo
		}

		convey.So(
			len(t.lfm.GetTopTracks(getAListOfArtists(sourceItems.Items), numberOfTopSongs).Items),
			assertion,
			equalValue...,
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
