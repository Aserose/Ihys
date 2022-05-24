package menu

import (
	"IhysBestowal/internal/datastruct"
)

type paginatedAudioItems struct {
	pageCount int
	from      string
	items     [][]datastruct.AudioItem
}

func paginateAudioItems(data datastruct.AudioItems, pageCapacity int) paginatedAudioItems {
	aip := paginatedAudioItems{
		pageCount: len(data.Items) / pageCapacity,
		from:      data.From,
		items:     make([][]datastruct.AudioItem, (len(data.Items)/pageCapacity)+1),
	}

	for i, j := 0, 0; i <= aip.pageCount; i, j = i+1, j+pageCapacity {
		var sss []datastruct.AudioItem

		if j+pageCapacity > len(data.Items) {
			sss = data.Items[j:]
		} else {
			sss = data.Items[j : j+pageCapacity]
		}

		aip.items[i] = sss
	}

	return aip
}

type paginatedPlaylistItems struct {
	pageCount int
	from      string
	items     [][]datastruct.PlaylistItem
}

func paginatePlaylistItems(data datastruct.PlaylistItems, pageCapacity int) paginatedPlaylistItems {
	ppi := paginatedPlaylistItems{
		pageCount: len(data.Items) / pageCapacity,
		from:      data.From,
		items:     make([][]datastruct.PlaylistItem, (len(data.Items)/pageCapacity)+1),
	}

	for i, j := 0, 0; i <= ppi.pageCount; i, j = i+1, j+pageCapacity {
		var sss []datastruct.PlaylistItem

		if j+pageCapacity > len(data.Items) {
			sss = data.Items[j:]
		} else {
			sss = data.Items[j : j+pageCapacity]
		}

		ppi.items[i] = sss
	}

	return ppi
}
