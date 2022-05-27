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
	lfm := newLfm(log)
	userId := int64(0)
	sourceItems := newSourceItems(
		newSong("Reliq", "gem"),
		newSong("Uniforms", "Serena"))

	convey.Convey("init", T, func() {

		convey.Convey("similar", func() {
			convey.So(
				lfm.GetSimiliarSongsFromLast(userId, sourceItems),
				convey.ShouldNotEqual,
				datastruct.AudioItems{})

		})

		convey.Convey("similiar100", func() {
			convey.So(
				len(lfm.GetSimiliarSongsFromLast100(userId, sourceItems).Items),
				convey.ShouldBeGreaterThan,
				3*len(sourceItems.Items))
		})

		convey.Convey("top", func() {
			getTopTracks := func(numberOfTopSongs int) {
				equalValue := numberOfTopSongs*len(sourceItems.Items)
				assertion := convey.ShouldEqual
				if numberOfTopSongs < 0 { equalValue = 0 }
				if numberOfTopSongs > 6 { assertion = convey.ShouldBeLessThanOrEqualTo }

				convey.So(
					len(lfm.GetTopTracks(getAListOfArtists(sourceItems.Items), numberOfTopSongs).Items),
					assertion, equalValue)
			}

			for _, num := range []int{1,3,5,-1} {
				getTopTracks(num)
			}

		})
	})
}

func newSong(artist, songTitle string) datastruct.AudioItem{
	return datastruct.AudioItem{
		Artist: artist,
		Title: songTitle,
	}
}

func newSourceItems(songs ... datastruct.AudioItem) datastruct.AudioItems{
	return datastruct.AudioItems{
		Items: songs,
	}
}

func getAListOfArtists(items []datastruct.AudioItem) []string {
	result := make([]string, len(items))
	for i, item := range items { result[i] = getArtist(item) }
	return result
}

func getArtist(item datastruct.AudioItem) string {
	return item.Artist
}


func newLfm(log customLogger.Logger) ILastFM {
	cfg := config.NewCfg(log)

	return  NewLastFM(
		log,
		config.NewCfg(log).LastFM,
		repository.NewRepository(log, cfg.Repository))
}