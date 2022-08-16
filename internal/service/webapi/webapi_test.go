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
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	testPostgres := map[string]string{
		"PSQL_USER":     "postgres",
		"PSQL_PASSWORD": "postgres",
		"PSQL_PORT":     "5432",
		"PSQL_HOST":     "localhost",
		"PSQL_NAME":     "postgres",
		"PSQL_SSLMODE":  "disable",
	}
	for k := range testPostgres {
		if err := os.Setenv(k, testPostgres[k]); err != nil {
			log.Print(err.Error())
		}
	}
	m.Run()
}

func TestWebapi(t *testing.T) {
	logs := customLogger.NewLogger()
	wa := newTestWebApi(logs)
	src := newSrc(
		newSong("Reliq", "gem"),
		newSong("Uniforms", "Serena"),
		newSong("Telepopmusik", "Close"))

	defer wa.Close()

	convey.Convey(" ", t, func() {

		convey.Convey("similar", func() { wa.similar(src) })
		convey.Convey("find", func() { wa.find() })

	})

}

type testWebApi struct {
	WebApi
}

func newTestWebApi(log customLogger.Logger) testWebApi {
	cfg := config.New(log)
	rep := repository.New(log, cfg.Repository)

	return testWebApi{
		WebApi: New(
			log,
			cfg.Service,
			rep,
			auth.New(log, cfg.Auth, rep)),
	}
}

func (t testWebApi) find() {
	convey.So(t.Find("does 214 it offend you we are").Title, convey.ShouldEqual, "We Are Rockstars")
}

func (t testWebApi) similar(src datastruct.Songs) {
	similar := func(perSource int) {
		equalValue := perSource * len(src.Songs)
		assertion := convey.ShouldBeGreaterThanOrEqualTo
		if perSource < 0 {
			equalValue = 0
		}
		if perSource > 10 {
			equalValue = 10
		}

		convey.So(len(t.WebApi.Similar(src, Opt{
			OnePerArtist: true,
			Ya: []yaMusic.Set{
				yaMusic.MaxPerSource(perSource),
			},
			Lf: []lastFm.Set{
				lastFm.MaxPerSource(perSource),
			},
		}).Songs), assertion, equalValue)
	}

	for _, num := range []int{2, 7, 3} {
		similar(num)
	}
}

func newSong(artist, title string) datastruct.Song {
	return datastruct.Song{
		Artist: artist,
		Title:  title,
	}
}

func newSrc(s ...datastruct.Song) datastruct.Songs {
	return datastruct.Songs{
		Songs: s,
	}
}
