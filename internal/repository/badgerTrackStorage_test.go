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
	db, err := newBadger(); if err != nil {
		logs.Panic(logs.CallInfoStr(), err.Error())
	}

	items := [17]datastruct.AudioItem{newItem("ok", "ok")}
	for i := 1; i < 17; i++ {
		items[i] = newItem(strconv.Itoa(i), strconv.Itoa(i))
	}

	tts := newTestTrackStorage(logs, db, newTestItems(items[:]...))

	convey.Convey("init", t, func() {

		convey.Convey("input-output", func() { tts.io() })
		convey.Convey("page handling", func() { tts.ioPage() })
		convey.Convey("delete", func() { tts.delete() })

	})

}

type testTrackStorage struct {
	storage badgerTrackStorage
	item    testItems
}

func newTestTrackStorage(log customLogger.Logger, db *badger.DB, items testItems) testTrackStorage {
	return testTrackStorage{
		storage: newBadgerTrackStorage(log, db),
		item:    items,
	}
}

func (t testTrackStorage) io() {
	convey.So(t.storage.GetItems(t.storage.Put(t.item.source, t.item.items), 0), convey.ShouldResemble, t.item.items.Items[:t.storage.GetPageCapacity()])
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

		convey.So(t.storage.GetPageCount(sourceAudio), assertion, equalValue)
	}

	for _, source := range []string{t.item.source.GetSourceAudio(t.item.items.From), "nonexistent 1", "2 nonexistent"} {
		getPageCount(source)
	}
}

func (t testTrackStorage) getItemsWithPageNum() {
	sourceAudio := t.item.source.GetSourceAudio(t.item.items.From)

	getItems := func(page int) {
		firstElementIdx := page * t.storage.GetPageCapacity()
		lastElementIdx := firstElementIdx + t.storage.GetPageCapacity()

		equalValue := []datastruct.AudioItem{}
		assertion := convey.ShouldResemble

		if len(t.item.items.Items) > firstElementIdx && firstElementIdx >= 0 {
			if lastElementIdx > len(t.item.items.Items) {
				equalValue = t.item.items.Items[firstElementIdx:]
			}
			if 0 < page || page < t.storage.GetPageCount(sourceAudio) {
				equalValue = t.item.items.Items[firstElementIdx:lastElementIdx]
			}
		}

		convey.So(t.storage.GetItems(sourceAudio, page), assertion, equalValue)
	}

	for _, selectedPage := range []int{1, 12, 0, -3} {
		getItems(selectedPage)
	}
}

func (t testTrackStorage) delete() {
	sourceAudio := t.item.source.GetSourceAudio(t.item.items.From)
	for _, assertion := range []convey.Assertion{convey.ShouldBeTrue, convey.ShouldBeFalse} {
		convey.So(t.storage.IsExist(sourceAudio), assertion)
		t.storage.Delete(sourceAudio)
	}
}

type testItems struct {
	source datastruct.AudioItem
	items  datastruct.AudioItems
}

func newTestItems(items ...datastruct.AudioItem) testItems {
	return testItems{
		source: items[0],
		items: datastruct.AudioItems{
			From:  "test",
			Items: items[1:],
		},
	}
}

func newItem(artist, title string) datastruct.AudioItem {
	return datastruct.AudioItem{
		Artist: artist,
		Title:  title,
	}
}
