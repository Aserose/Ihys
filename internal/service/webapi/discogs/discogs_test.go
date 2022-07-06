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

	convey.Convey(`init`, t, func() {

		convey.Convey(`artist website`, func() { d.getArtistWebsite(artists) })
		convey.Convey(`song info`, func() { d.getSongInfo(songs) })

	})
}

type testDiscogs struct {
	enquirer
	collater
}

func newTestDiscogs(log customLogger.Logger) testDiscogs {
	return testDiscogs{
		collater: newCollater(),
		enquirer: newEnquirer(log, config.NewCfg(log).Service.Discogs),
	}
}

func (t testDiscogs) getSongInfo(songs []datastruct.AudioItem) {
	getSong := func(song datastruct.AudioItem) {
		equalValue := datastruct.AudioInfo{}
		assertion := convey.ShouldNotResemble
		if song.Artist == `` {
			assertion = convey.ShouldResemble
		}

		convey.So(t.enquirer.getSongInfo(song), assertion, equalValue)
	}

	for _, song := range songs {
		getSong(song)
	}
}

func (t testDiscogs) getArtistWebsite(artists []string) {
	getWebsite := func(artist string) {
		assertion := convey.ShouldNotBeEmpty

		if artist == `` {
			assertion = convey.ShouldBeEmpty
		}

		convey.So(t.enquirer.getWebsites(artist, typeArtist), assertion)
	}

	for _, artist := range artists {
		getWebsite(artist)
	}
}

func newTestAudios(audios ...datastruct.AudioItem) []datastruct.AudioItem {
	return audios
}

func newTestAudio(artist, title string) datastruct.AudioItem {
	return datastruct.AudioItem{
		Artist: artist,
		Title:  title,
	}
}
