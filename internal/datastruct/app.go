package datastruct

import "strings"

type Playlists struct {
	From      string
	Playlists []Playlist
}

type Playlist struct {
	ID      int
	OwnerId int
	Title   string
}

type Songs struct {
	From  string
	Songs []Song
}

func (a Songs) WithFrom(i int) string {
	return a.Songs[i].WithFrom(a.From)
}

type Song struct {
	Artist string `json:"artist"`
	Title  string `json:"title"`
	Url    string `json:"url"`
}

func (a Song) FirstArtist() string {
	if strings.Contains(a.Artist, `, `) {
		return strings.Split(a.Artist, `, `)[0]
	}
	return a.Artist
}

func (a Song) NewSongs(src string) Songs {
	leftSepar, rightSepar := a.Separators()
	song := strings.Split(src, ` - `)

	if !strings.Contains(src, leftSepar) {
		return Songs{
			Songs: []Song{
				{
					Artist: song[0],
					Title:  strings.Split(song[1], "\n\n")[0],
				},
			},
		}
	}

	s := strings.Split(song[1], leftSepar)

	return Songs{
		From: strings.Replace(s[1], rightSepar, ``, 1),
		Songs: []Song{
			{
				Artist: song[0],
				Title:  s[0],
			},
		},
	}
}

func (a Song) Separators() (left string, right string) {
	return ` «(`, `)»`
}

func (a Song) WithFrom(from string) string {
	leftSep, rightSep := a.Separators()
	return a.Artist + ` - ` + a.Title + leftSep + from + rightSep
}

func (a Song) WithoutFrom() string {
	return a.Artist + ` - ` + a.Title
}

type ExecParam struct {
	ChatId int64
	MsgId  int
}

type SongInfo struct {
	Label       string
	Genres      []string
	Country     string
	ReleaseDate string
}

func (a SongInfo) Genre() string {
	return strings.Join(a.Genres, `, `)
}
