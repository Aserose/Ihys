package lastFm

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/repository"
	"IhysBestowal/pkg/customLogger"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

//func TestMain(m *testing.M) {
//	testPostgres := map[string]string{
//		"PSQL_USER":     "postgres",
//		"PSQL_PASSWORD": "postgres",
//		"PSQL_PORT":     "5432",
//		"PSQL_HOST":     "localhost",
//		"PSQL_NAME":     "postgres",
//		"PSQL_SSLMODE":  "disable",
//	}
//
//	for k, v := range testPostgres {
//		if err := os.Setenv(k, v); err != nil {
//			log.Print(err.Error())
//		}
//	}
//
//	m.Run()
//}

func TestLastFm(T *testing.T) {
	logs := customLogger.NewLogger()
	lfm := newTestLfm(logs)
	uid := int64(0)
	src := newSrc(
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

func (t testLfm) similar(uid int64, src datastruct.Set) {
	sim := func(amountPerSource int) {
		equalValue := amountPerSource * len(src.Song)
		assertion := convey.ShouldEqual
		if amountPerSource < 0 {
			equalValue = 0
		}
		if amountPerSource > 6 {
			equalValue = amountPerSource * 6
			assertion = convey.ShouldBeGreaterThanOrEqualTo
		}

		convey.So(len(t.Similar(uid, src, MaxPerSource(amountPerSource)).Song),
			assertion, equalValue)
	}

	for _, num := range []int{6, 1, 0, -4, 2} {
		sim(num)
	}

}

func (t testLfm) top(src datastruct.Set) {
	get := func(numSongs int) {
		equalValue := numSongs * len(src.Song)
		assertion := convey.ShouldEqual
		if numSongs < 0 {
			equalValue = 0
		}
		if numSongs > 6 {
			equalValue = numSongs * 6
			assertion = convey.ShouldBeGreaterThanOrEqualTo
		}

		convey.So(len(t.Top(artists(src.Song), numSongs).Song), assertion, equalValue)
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

func newSrc(songs ...datastruct.Song) datastruct.Set {
	return datastruct.Set{
		Song: songs,
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

	return New(log, config.New(log).LastFM, repository.New(log, cfg.Repository))
}
