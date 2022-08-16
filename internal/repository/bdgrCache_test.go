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

func TestBdgr(t *testing.T) {
	logs := customLogger.NewLogger()
	db := newBdgr(logs)
	defer db.Close()

	items := [17]datastruct.Song{newItem("ok", "ok")}
	for i := 1; i < 17; i++ {
		items[i] = newItem(strconv.Itoa(i), strconv.Itoa(i))
	}

	c := newTestCache(logs, db, newTestItems(items[:]...))

	convey.Convey(" ", t, func() {

		convey.Convey("input-output", func() { c.io() })
		convey.Convey("page handling", func() { c.ioPage() })
		convey.Convey("drop", func() { c.drop() })
		convey.Convey("drop all", func() { c.dropAll() })

	})

}

type testCache struct {
	cache bdgrCache
	val   testItems
}

func newTestCache(log customLogger.Logger, db *badger.DB, items testItems) testCache {
	return testCache{
		cache: newBdgrCache(log, db),
		val:   items,
	}
}

func (t testCache) io() {
	convey.So(t.cache.Get(t.cache.Put(t.val.src, t.val.items), 0), convey.ShouldResemble, t.val.items.Songs[:t.cache.PageCapacity()])
}

func (t testCache) ioPage() {
	t.getPageCount()
	t.itemsWithPageNum()
}

func (t testCache) getPageCount() {
	getPageCount := func(sourceAudio string) {
		equalValue := -1
		assertion := convey.ShouldNotEqual

		if strings.Contains(sourceAudio, "nonexistent") {
			assertion = convey.ShouldEqual
		}

		convey.So(t.cache.PageCount(sourceAudio), assertion, equalValue)
	}

	for _, source := range []string{t.val.src.WithFrom(t.val.items.From), "nonexistent 1", "2 nonexistent"} {
		getPageCount(source)
	}
}

func (t testCache) itemsWithPageNum() {
	src := t.val.src.WithFrom(t.val.items.From)

	get := func(page int) {
		first := page * t.cache.PageCapacity()
		last := first + t.cache.PageCapacity()

		equalValue := []datastruct.Song{}
		assertion := convey.ShouldResemble

		if page >= 0 {
			if len(t.val.items.Songs) > first && first >= 0 {
				if last > len(t.val.items.Songs) {
					equalValue = t.val.items.Songs[first:]
				}
				if 0 < page || page < t.cache.PageCount(src) {
					equalValue = t.val.items.Songs[first:last]
				}
			} else {
				equalValue = t.val.items.Songs[t.cache.PageCount(src)*t.cache.pageCapacity:]
			}
		}
		convey.So(t.cache.Get(src, page), assertion, equalValue)
	}

	for _, pageNum := range []int{1, 12, 0, -3} {
		get(pageNum)
	}
}

func (t testCache) dropAll() {
	t.io()
	src := t.val.src.WithFrom(t.val.items.From)
	for _, assertion := range []convey.Assertion{convey.ShouldBeTrue, convey.ShouldBeFalse} {
		convey.So(t.cache.IsExist(src), assertion)
		t.cache.dropAll()
	}
}

func (t testCache) drop() {
	sourceItems := t.val.src.WithFrom(t.val.items.From)
	for _, assertion := range []convey.Assertion{convey.ShouldBeTrue, convey.ShouldBeFalse} {
		convey.So(t.cache.IsExist(sourceItems), assertion)
		t.cache.drop(sourceItems)
	}
}

type testItems struct {
	src   datastruct.Song
	items datastruct.Songs
}

func newTestItems(items ...datastruct.Song) testItems {
	return testItems{
		src: items[0],
		items: datastruct.Songs{
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
