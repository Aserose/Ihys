package datastruct

import "strings"

type PlaylistItems struct {
	From  string
	Items []PlaylistItem
}

type PlaylistItem struct {
	ID      int
	OwnerId int
	Title   string
}

type AudioItems struct {
	From  string
	Items []AudioItem
}

type AudioItem struct {
	Artist string `json:"artist"`
	Title  string `json:"title"`
	Url    string `json:"url"`
}

func (a AudioItem) GetFirstArtist() string {
	if strings.Contains(a.Artist, `, `) {
		return strings.Split(a.Artist, `, `)[0]
	}
	return a.Artist
}

func (a AudioItem) GetSourceAudio(sourceFrom string) string {
	return a.Artist + ` - ` + a.Title + ` (` + sourceFrom + `)`
}

func (a AudioItem) GetAudio() string {
	return a.Artist + ` - ` + a.Title
}

type Datum struct {
	UserId int64
	ChatID int64
	MsgId  int
	Data   string
}

type ExecParam struct {
	ChatId int64
	MsgId  int
}

type AudioInfo struct {
	Label   string
	Genres  []string
	Country string
	Year    string
}

func (a AudioInfo) GetGenresString() string {
	return strings.Join(a.Genres, `, `)
}
