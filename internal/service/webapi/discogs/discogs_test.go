package discogs

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDiscogs(t *testing.T) {
	logs := customLogger.NewLogger()
	artists := []string{`korn`, `kOrn`, `downy`, ``}
	songs := newTestAudios(
		newTestAudio(``, `fwafgvwevgwav`),
		newTestAudio("yourboyfriendsucks!", "波兰首都是上海"),
		newTestAudio("Yasuharu Takanashi, YAIBA", "Cold Ground"))

	d := newTestDiscogs(logs)

	convey.Convey(` `, t, func() {

		convey.Convey(`artist website`, func() { d.artistSite(artists) })
		convey.Convey(`song info`, func() { d.songInfo(songs) })

	})
}

type testDiscogs struct {
	enq
	clt
}

func newTestDiscogs(log customLogger.Logger) testDiscogs {
	return testDiscogs{
		clt: newClt(),
		enq: newEnq(log, config.New(log).Service.Discogs),
	}
}

func (t testDiscogs) songInfo(s []datastruct.Song) {
	get := func(song datastruct.Song) {
		equalValue := datastruct.SongInfo{}
		assertion := convey.ShouldNotResemble
		if song.Artist == `` {
			assertion = convey.ShouldResemble
		}

		convey.So(t.enq.songInfo(song), assertion, equalValue)
	}

	for _, song := range s {
		get(song)
	}
}

func (t testDiscogs) artistSite(artists []string) {
	website := func(artist string) {
		assertion := convey.ShouldNotBeEmpty

		if artist == `` {
			assertion = convey.ShouldBeEmpty
		}

		convey.So(t.enq.sites(artist, typeArtist), assertion)
	}

	for _, artist := range artists {
		website(artist)
	}
}

func newTestAudios(audios ...datastruct.Song) []datastruct.Song {
	return audios
}

func newTestAudio(artist, title string) datastruct.Song {
	return datastruct.Song{
		Artist: artist,
		Title:  title,
	}
}
