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

func (a AudioItems) GetSourceAudio(elemNum int) string {
	return a.Items[elemNum].GetSourceAudio(a.From)
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

func (a AudioItem) GetSeparators() (left string, right string) {
	return ` «(`, `)»`
}

func (a AudioItem) GetSourceAudio(sourceFrom string) string {
	leftSep, rightSep := a.GetSeparators()
	return a.Artist + ` - ` + a.Title + leftSep + sourceFrom + rightSep
}

func (a AudioItem) GetAudio() string {
	return a.Artist + ` - ` + a.Title
}

type ExecParam struct {
	ChatId int64
	MsgId  int
}

type AudioInfo struct {
	Label       string
	Genres      []string
	Country     string
	ReleaseDate string
}

func (a AudioInfo) GetGenresString() string {
	return strings.Join(a.Genres, `, `)
}
