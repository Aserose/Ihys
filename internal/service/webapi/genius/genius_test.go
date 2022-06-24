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
	gen := NewGenius(log, config.NewCfg(log).Service.Genius)
	testItem := datastruct.AudioItem{
		Artist: "Violet Cold",
		Title:  "Anomie",
	}

	convey.Convey(`init`, t, func() {

		convey.So(gen.GetLyricsURL(testItem), convey.ShouldNotBeEmpty)

	})

}
