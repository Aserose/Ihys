package genius

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGenius(t *testing.T) {
	log := customLogger.NewLogger()
	gen := New(log, config.New(log).Service.Genius)
	testItem := datastruct.Song{
		Artist: "Violet Cold",
		Title:  "Anomie",
	}

	convey.Convey(`init`, t, func() {

		convey.So(gen.LyricsURL(testItem), convey.ShouldNotBeEmpty)

	})

}
