package repository

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
	"github.com/dgraph-io/badger/v3"
	"github.com/smartystreets/goconvey/convey"
	"strconv"
	"strings"
	"testing"
)

func TestBadger(t *testing.T) {
	logs := customLogger.NewLogger()
	db := newBadger(logs)
	defer db.Close()

	items := [17]datastruct.Song{newItem("ok", "ok")}
	for i := 1; i < 17; i++ {
		items[i] = newItem(strconv.Itoa(i), strconv.Itoa(i))
	}

	tts := newTestTrackStorage(logs, db, newTestItems(items[:]...))

	convey.Convey("init", t, func() {

		convey.Convey("input-output", func() { tts.io() })
		convey.Convey("page handling", func() { tts.ioPage() })
		convey.Convey("drop", func() { tts.drop() })
		convey.Convey("drop all", func() { tts.dropAll() })

	})

}

type testTrackStorage struct {
	storage bdgrCache
	item    testItems
}

func newTestTrackStorage(log customLogger.Logger, db *badger.DB, items testItems) testTrackStorage {
	return testTrackStorage{
		storage: newBdgrCache(log, db),
		item:    items,
	}
}

func (t testTrackStorage) io() {
	convey.So(t.storage.Get(t.storage.Put(t.item.source, t.item.subItems), 0), convey.ShouldResemble, t.item.subItems.Songs[:t.storage.PageCapacity()])
}

func (t testTrackStorage) ioPage() {
	t.getPageCount()
	t.getItemsWithPageNum()
}

func (t testTrackStorage) getPageCount() {
	getPageCount := func(sourceAudio string) {
		equalValue := -1
		assertion := convey.ShouldNotEqual

		if strings.Contains(sourceAudio, "nonexistent") {
			assertion = convey.ShouldEqual
		}

		convey.So(t.storage.PageCount(sourceAudio), assertion, equalValue)
	}

	for _, source := range []string{t.item.source.WithFrom(t.item.subItems.From), "nonexistent 1", "2 nonexistent"} {
		getPageCount(source)
	}
}

func (t testTrackStorage) getItemsWithPageNum() {
	sourceItems := t.item.source.WithFrom(t.item.subItems.From)

	getItems := func(page int) {
		firstElementIdx := page * t.storage.PageCapacity()
		lastElementIdx := firstElementIdx + t.storage.PageCapacity()

		equalValue := []datastruct.Song{}
		assertion := convey.ShouldResemble

		if page >= 0 {
			if len(t.item.subItems.Songs) > firstElementIdx && firstElementIdx >= 0 {
				if lastElementIdx > len(t.item.subItems.Songs) {
					equalValue = t.item.subItems.Songs[firstElementIdx:]
				}
				if 0 < page || page < t.storage.PageCount(sourceItems) {
					equalValue = t.item.subItems.Songs[firstElementIdx:lastElementIdx]
				}
			} else {
				equalValue = t.item.subItems.Songs[t.storage.PageCount(sourceItems)*t.storage.pageCapacity:]
			}
		}
		convey.So(t.storage.Get(sourceItems, page), assertion, equalValue)
	}

	for _, selectedPage := range []int{1, 12, 0, -3} {
		getItems(selectedPage)
	}
}

func (t testTrackStorage) dropAll() {
	t.io()
	sourceItems := t.item.source.WithFrom(t.item.subItems.From)
	for _, assertion := range []convey.Assertion{convey.ShouldBeTrue, convey.ShouldBeFalse} {
		convey.So(t.storage.IsExist(sourceItems), assertion)
		t.storage.dropAll()
	}
}

func (t testTrackStorage) drop() {
	sourceItems := t.item.source.WithFrom(t.item.subItems.From)
	for _, assertion := range []convey.Assertion{convey.ShouldBeTrue, convey.ShouldBeFalse} {
		convey.So(t.storage.IsExist(sourceItems), assertion)
		t.storage.drop(sourceItems)
	}
}

type testItems struct {
	source   datastruct.Song
	subItems datastruct.Songs
}

func newTestItems(items ...datastruct.Song) testItems {
	return testItems{
		source: items[0],
		subItems: datastruct.Songs{
			From:  "test",
			Songs: items[1:],
		},
	}
}

func newItem(artist, title string) datastruct.Song {
	return datastruct.Song{
		Artist: artist,
		Title:  title,
	}
}
