package pgntr

import "IhysBestowal/internal/datastruct"

type Songs struct {
	PageCount int                 `json:"page_count"`
	AddDate   int                 `json:"add_date"`
	Items     [][]datastruct.Song `json:"items"`
}

func NewSongs(data datastruct.Songs, pageCapacity int) Songs {
	s := Songs{
		PageCount: len(data.Songs) / pageCapacity,
		Items:     make([][]datastruct.Song, (len(data.Songs)/pageCapacity)+1),
	}

	for i, j := 0, 0; i <= s.PageCount; i, j = i+1, j+pageCapacity {
		var items []datastruct.Song

		if j+pageCapacity > len(data.Songs) {
			items = data.Songs[j:]
		} else {
			items = data.Songs[j : j+pageCapacity]
		}

		s.Items[i] = items
	}

	return s
}
