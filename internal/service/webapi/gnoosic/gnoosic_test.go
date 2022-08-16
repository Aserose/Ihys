package gnoosic

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGnoosic(t *testing.T) {
	convey.Convey(`init`, t, func() {
		convey.So(New().RandomArtist(), convey.ShouldNotEqual, `/`)
	})
}
