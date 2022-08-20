package pgntr

// TODO
//type Playlists struct {
//	PageCount int
//	From      string
//	Items     [][]datastruct.Playlist
//}
//
//func NewPlaylists(data datastruct.Playlists, pageCapacity int) Playlists {
//	p := Playlists{
//		PageCount: len(data.Playlists) / pageCapacity,
//		From:      data.From,
//		Items:     make([][]datastruct.Playlist, (len(data.Playlists)/pageCapacity)+1),
//	}
//
//	for i, j := 0, 0; i <= p.PageCount; i, j = i+1, j+pageCapacity {
//		var items []datastruct.Playlist
//
//		if j+pageCapacity > len(data.Playlists) {
//			items = data.Playlists[j:]
//		} else {
//			items = data.Playlists[j : j+pageCapacity]
//		}
//
//		p.Items[i] = items
//	}
//
//	return p
//}
