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

	items := [17]datastruct.AudioItem{newItem("ok", "ok")}
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
	convey.So(t.storage.GetItems(t.storage.Put(t.item.source, t.item.subItems), 0), convey.ShouldResemble, t.item.subItems.Items[:t.storage.GetPageCapacity()])
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

	for _, source := range []string{t.item.source.GetSourceAudio(t.item.subItems.From), "nonexistent 1", "2 nonexistent"} {
		getPageCount(source)
	}
}

func (t testTrackStorage) getItemsWithPageNum() {
	sourceItems := t.item.source.GetSourceAudio(t.item.subItems.From)

	getItems := func(page int) {
		firstElementIdx := page * t.storage.GetPageCapacity()
		lastElementIdx := firstElementIdx + t.storage.GetPageCapacity()

		equalValue := []datastruct.AudioItem{}
		assertion := convey.ShouldResemble

		if page >= 0 {
			if len(t.item.subItems.Items) > firstElementIdx && firstElementIdx >= 0 {
				if lastElementIdx > len(t.item.subItems.Items) {
					equalValue = t.item.subItems.Items[firstElementIdx:]
				}
				if 0 < page || page < t.storage.GetPageCount(sourceItems) {
					equalValue = t.item.subItems.Items[firstElementIdx:lastElementIdx]
				}
			} else {
				equalValue = t.item.subItems.Items[t.storage.GetPageCount(sourceItems)*t.storage.pageCapacity:]
			}
		}
		convey.So(t.storage.GetItems(sourceItems, page), assertion, equalValue)
	}

	for _, selectedPage := range []int{1, 12, 0, -3} {
		getItems(selectedPage)
	}
}

func (t testTrackStorage) dropAll() {
	t.io()
	sourceItems := t.item.source.GetSourceAudio(t.item.subItems.From)
	for _, assertion := range []convey.Assertion{convey.ShouldBeTrue, convey.ShouldBeFalse} {
		convey.So(t.storage.IsExist(sourceItems), assertion)
		t.storage.dropAll()
	}
}

func (t testTrackStorage) drop() {
	sourceItems := t.item.source.GetSourceAudio(t.item.subItems.From)
	for _, assertion := range []convey.Assertion{convey.ShouldBeTrue, convey.ShouldBeFalse} {
		convey.So(t.storage.IsExist(sourceItems), assertion)
		t.storage.drop(sourceItems)
	}
}

type testItems struct {
	source   datastruct.AudioItem
	subItems datastruct.AudioItems
}

func newTestItems(items ...datastruct.AudioItem) testItems {
	return testItems{
		source: items[0],
		subItems: datastruct.AudioItems{
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
