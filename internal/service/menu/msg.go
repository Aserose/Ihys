package menu

import (
	"IhysBestowal/internal/datastruct"
	"github.com/biter777/countries"
	"strings"
)

const (
	rightArrow = "->"
	leftArrow  = "<-"
	vLeft      = "vLeft"
	vRight     = "vRight"
	vSelect    = "vSelect"
	pageNum    = "pageNum"
	spc        = ` `

	backTxt = "back"
	backClb = "back"
	emp     = ``

	emojiMovieCamera  = " \xF0\x9F\x8E\xA5 "
	emojiLink         = " \xF0\x9F\x94\x97 "
	emojiHourglass    = " \xE2\x8C\x9B "
	emojiBlackNim     = " \xE2\x9C\x92 "
	emojiPageWithCurl = " \xF0\x9F\x93\x83 "

	separator      = ` | `
	msgYouTube     = separator + emojiMovieCamera + `[YouTube]`
	msgWebsite     = separator + emojiLink + `[Website]`
	msgLyrics      = separator + emojiPageWithCurl + `[Lyrics]`
	msgLoadingBase = emojiHourglass + `Un momento! It's uploading.`

	idt    = "\n"
	dblIdt = "\n\n"
)

var msgLoading = [...]string{
	formatQuote(`Patience attracts happiness; it brings near that which is far.`, `Swahili Proverb`),
	formatQuote(`Our patience will achieve more than our force.`, `Edmund Burke`),
	formatQuote(`Learning patience can be a difficult experience, but once conquered, you will find life is easier.`, `Catherine Pulsifer`),
	formatQuote(`Patience is the best remedy for every trouble.`, `Plautus`),
	formatQuote(`Trees that are slow to grow bear the best fruit`, `Moliere`),
}

func formatSong(song datastruct.Song) string {
	return "\n" + "[" + song.WithoutFrom() + "]" + "(" + song.Url + ")"
}

func formatQuote(quote, author string) string {
	return msgLoadingBase + "\n\n" + "“" + quote + "“ \n - " + author + emojiBlackNim
}

func formatInfo(info datastruct.SongInfo) string {
	var flg string
	if country := countries.ByName(info.Country); country.String() != countries.UnknownMsg {
		flg = country.Emoji()
	}

	return build(`Label: `, info.Label, ` < `, info.Country, `  `, flg, ` > `,
		idt, `Release: `, info.ReleaseDate,
		idt, `Genre: `, info.Genre(), dblIdt)
}

func formatVideoURL(url string) string  { return msgYouTube + formatURL(url) }
func formatLyricsURL(url string) string { return msgLyrics + formatURL(url) }
func formatWebsite(url string) string   { return msgWebsite + formatURL(url) }
func formatURL(url string) string       { return `(` + url + `)` }

func build(str ...string) string {
	b := strings.Builder{}

	for _, s := range str {
		b.WriteString(s)
	}

	return b.String()
}
