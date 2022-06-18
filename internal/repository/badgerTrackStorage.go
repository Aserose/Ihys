package repository

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
	"github.com/dgraph-io/badger/v3"
	"github.com/goccy/go-json"
	"github.com/tidwall/gjson"
)

type badgerTrackStorage struct {
	db           *badger.DB
	pageCapacity int
	log          customLogger.Logger
}

func newBadgerTrackStorage(log customLogger.Logger, db *badger.DB) badgerTrackStorage {
	return badgerTrackStorage{
		db:           db,
		pageCapacity: 7,
		log:          log,
	}
}

func (b badgerTrackStorage) Put(sourceAudio datastruct.AudioItem, similar datastruct.AudioItems) string {
	data, _ := json.Marshal(paginateAudioItems(similar, b.pageCapacity))
	key := sourceAudio.GetSourceAudio(similar.From)

	if err := b.db.Update(func(txn *badger.Txn) error { return txn.Set([]byte(key), data) }); err != nil {
		b.log.Panic(b.log.CallInfoStr(), err.Error())
		return ``
	}

	return key
}

func (b badgerTrackStorage) GetPageCapacity() int {
	return b.pageCapacity
}

func (b badgerTrackStorage) GetPageCount(sourceAudio string) int {
	var count int
	if err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(sourceAudio))
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

func (b badgerTrackStorage) IsExist(sourceAudio string) bool {
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(sourceAudio))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error { return nil })
	})

	return err == nil
}

func (b badgerTrackStorage) GetItems(sourceAudio string, page int) []datastruct.AudioItem {
	if !b.IsExist(sourceAudio) || 0 > page || page > b.GetPageCount(sourceAudio) {
		return []datastruct.AudioItem{}
	}

	pai := paginatedAudioItems{}

	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(sourceAudio))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error { return json.Unmarshal(val, &pai) })
	})
	if err != nil {
		return []datastruct.AudioItem{}
	}

	return pai.Items[page]
}

func (b badgerTrackStorage) Delete(sourceAudio string) {
	if b.IsExist(sourceAudio) {
		if err := b.db.Update(func(txn *badger.Txn) error { return txn.Delete([]byte(sourceAudio)) }); err != nil {
			b.log.Error(b.log.CallInfoStr(), err.Error())
		}
	}
}

type paginatedAudioItems struct {
	PageCount int                      `json:"page_count"`
	Items     [][]datastruct.AudioItem `json:"items"`
}

func paginateAudioItems(data datastruct.AudioItems, pageCapacity int) paginatedAudioItems {
	pai := paginatedAudioItems{
		PageCount: len(data.Items) / pageCapacity,
		Items:     make([][]datastruct.AudioItem, (len(data.Items)/pageCapacity)+1),
	}

	for i, j := 0, 0; i <= pai.PageCount; i, j = i+1, j+pageCapacity {
		var items []datastruct.AudioItem

		if j+pageCapacity > len(data.Items) {
			items = data.Items[j:]
		} else {
			items = data.Items[j : j+pageCapacity]
		}

		pai.Items[i] = items
	}

	return pai
}

type paginatedPlaylistItems struct {
	PageCount int
	From      string
	Items     [][]datastruct.PlaylistItem
}

func paginatePlaylistItems(data datastruct.PlaylistItems, pageCapacity int) paginatedPlaylistItems {
	ppi := paginatedPlaylistItems{
		PageCount: len(data.Items) / pageCapacity,
		From:      data.From,
		Items:     make([][]datastruct.PlaylistItem, (len(data.Items)/pageCapacity)+1),
	}

	for i, j := 0, 0; i <= ppi.PageCount; i, j = i+1, j+pageCapacity {
		var items []datastruct.PlaylistItem

		if j+pageCapacity > len(data.Items) {
			items = data.Items[j:]
		} else {
			items = data.Items[j : j+pageCapacity]
		}

		ppi.Items[i] = items
	}

	return ppi
}
