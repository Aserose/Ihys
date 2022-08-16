package lastFm

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/repository"
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

func TestLastFm(T *testing.T) {
	logs := customLogger.NewLogger()
	lfm := newTestLfm(logs)
	uid := int64(0)
	src := newSourceItems(
		newSong("Reliq", "gem"),
		newSong("Uniforms", "Serena"),
		newSong("Telepopmusik", "Close"))

	convey.Convey(" ", T, func() {

		convey.Convey("similar", func() { lfm.similar(uid, src) })
		convey.Convey("top", func() { lfm.top(src) })

	})
}

type testLfm struct {
	LastFM
}

func newTestLfm(log customLogger.Logger) testLfm {
	return testLfm{
		LastFM: newLfm(log),
	}
}

func (t testLfm) similar(uid int64, src datastruct.Songs) {
	sim := func(amountPerSource int) {
		equalValue := amountPerSource * len(src.Songs)
		assertion := convey.ShouldEqual
		if amountPerSource < 0 {
			equalValue = 0
		}
		if amountPerSource > 6 {
			equalValue = amountPerSource * 6
			assertion = convey.ShouldBeGreaterThanOrEqualTo
		}

		convey.So(len(t.Similar(uid, src, MaxPerSource(amountPerSource)).Songs),
			assertion, equalValue)
	}

	for _, num := range []int{6, 1, 0, -4, 2} {
		sim(num)
	}

}

func (t testLfm) top(src datastruct.Songs) {
	get := func(numSongs int) {
		equalValue := numSongs * len(src.Songs)
		assertion := convey.ShouldEqual
		if numSongs < 0 {
			equalValue = 0
		}
		if numSongs > 6 {
			equalValue = numSongs * 6
			assertion = convey.ShouldBeGreaterThanOrEqualTo
		}

		convey.So(
			len(t.Top(artists(src.Songs), numSongs).Songs),
			assertion, equalValue,
		)
	}

	for _, num := range []int{1, 3, 5, -1} {
		get(num)
	}
}

func newSong(artist, title string) datastruct.Song {
	return datastruct.Song{
		Artist: artist,
		Title:  title,
	}
}

func newSourceItems(songs ...datastruct.Song) datastruct.Songs {
	return datastruct.Songs{
		Songs: songs,
	}
}

func artists(s []datastruct.Song) []string {
	result := make([]string, len(s))
	for i, item := range s {
		result[i] = artist(item)
	}
	return result
}

func artist(s datastruct.Song) string {
	return s.Artist
}

func newLfm(log customLogger.Logger) LastFM {
	cfg := config.New(log)

	return New(
		log,
		config.New(log).LastFM,
		repository.New(log, cfg.Repository))
}
