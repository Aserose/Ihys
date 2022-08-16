package repository

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/repository/pgntr"
	"IhysBestowal/pkg/customLogger"
	"github.com/dgraph-io/badger/v3"
	"github.com/goccy/go-json"
	"github.com/tidwall/gjson"
)

type bdgrCache struct {
	db           *badger.DB
	pageCapacity int
	maxSize      float32
	log          customLogger.Logger
}

func newBdgrCache(log customLogger.Logger, db *badger.DB) bdgrCache {
	return bdgrCache{
		db:           db,
		pageCapacity: 7,
		maxSize:      10.0,
		log:          log,
	}
}

func (b bdgrCache) Put(src datastruct.Song, similar datastruct.Songs) string {
	if b.size() > b.maxSize {
		b.dropAll()
	}
	data, _ := json.Marshal(pgntr.NewSongs(similar, b.pageCapacity))
	key := src.WithFrom(similar.From)

	if err := b.db.Update(func(txn *badger.Txn) error { return txn.Set([]byte(key), data) }); err != nil {
		b.log.Panic(b.log.CallInfoStr(), err.Error())
		return ``
	}

	return key
}

func (b bdgrCache) size() float32 {
	_, vlog := b.db.Size()
	return float32(vlog) / (1 << 30)
}

func (b bdgrCache) PageCapacity() int {
	return b.pageCapacity
}

func (b bdgrCache) PageCount(src string) int {
	var count int
	if err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(src))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			count = int(gjson.GetBytes(val, "page_count").Int())
			return nil
		})
	}); err != nil {
		b.log.Error(b.log.CallInfoStr(), err.Error())
		return -1
	}

	return count
}

func (b bdgrCache) IsExist(src string) bool {
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(src))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error { return nil })
	})

	return err == nil
}

func (b bdgrCache) Get(src string, page int) []datastruct.Song {
	if !b.IsExist(src) || 0 > page {
		return []datastruct.Song{}
	}

	res := pgntr.Songs{}

	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(src))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error { return json.Unmarshal(val, &res) })
	})
	if err != nil {
		return []datastruct.Song{}
	}

	if page > res.PageCount {
		return res.Items[len(res.Items)-1]
	}

	return res.Items[page]
}

func (b bdgrCache) dropAll() {
	if err := b.db.DropAll(); err != nil {
		b.log.Error(b.log.CallInfoStr(), err.Error())
	}
}

func (b bdgrCache) drop(src string) {
	if b.IsExist(src) {
		if err := b.db.Update(func(txn *badger.Txn) error { return txn.Delete([]byte(src)) }); err != nil {
			b.log.Error(b.log.CallInfoStr(), err.Error())
		}
	}
}
