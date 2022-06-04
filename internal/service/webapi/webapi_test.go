package webapi

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/repository"
	"IhysBestowal/internal/service/auth"
	"IhysBestowal/internal/service/webapi/lastFm"
	"IhysBestowal/internal/service/webapi/yaMusic"
	"IhysBestowal/pkg/customLogger"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestWebapi(t *testing.T) {
	logs := customLogger.NewLogger()
	webapi := newTestWebApi(logs)
	sourceItems := newSourceItems(
		newSong("Reliq", "gem"),
		newSong("Uniforms", "Serena"),
		newSong("Telepopmusik", "Close"))

	convey.Convey("init", t, func() {

		convey.Convey("get similar", func() { webapi.getSimilar(sourceItems) })
		convey.Convey("song search", func() { webapi.songSearch() })

	})

}

type testWebApi struct {
	WebApiService
}

func newTestWebApi(log customLogger.Logger) testWebApi {
	cfg := config.NewCfg(log)
	rep := repository.NewRepository(log, cfg.Repository)

	return testWebApi{
		WebApiService: NewWebApiService(
			log,
			cfg.Service,
			rep,
			auth.NewAuthService(log, cfg.Auth, rep)),
	}
}

func (t testWebApi) songSearch() {
	convey.So(t.GetAudio("does 214 it offend you we are").Title, convey.ShouldEqual, "We Are Rockstars")
}

func (t testWebApi) getSimilar(source datastruct.AudioItems) {
	getSimiliar := func(amountPerSource int) {
		equalValue := amountPerSource * len(source.Items)
		assetion := convey.ShouldBeGreaterThanOrEqualTo
		if amountPerSource < 0 {
			equalValue = 0
		}
		if amountPerSource > 10 {
			equalValue = 10
		}

		convey.So(len(t.WebApiService.GetSimilar(source, Opt{
			oneAudioPerArtist: true,
			ya: []yaMusic.ProcessingOptions{
				yaMusic.SetMaxAudioAmountPerSource(amountPerSource),
			},
			lf: []lastFm.ProcessingOptions{
				lastFm.SetMaxAudioAmountPerSource(amountPerSource),
			},
		}).Items), assetion, equalValue)
	}

	for _, num := range []int{2, 7, 3} {
		getSimiliar(num)
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
