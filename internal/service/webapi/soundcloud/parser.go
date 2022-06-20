package soundcloud

import (
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/tidwall/gjson"
	"github.com/v-grabko1999/go-html2json"
	"strings"
	"time"
	"unicode"
)

const (
	pathURLPathname   = "#content > div > div > div.l-main > div > div > div > ul > li:nth-child(%d) > div > div > div > div.sound__artwork.sc-mr-1x > a"
	pathFlexContainer = "#content > div > div > div.l-main.g-main-scroll-area > div > div > div.relatedList__list > div"
	pathItemList      = "#content > div > div > div.l-main.g-main-scroll-area > div > div > div.relatedList__list > ul"
	pathItem          = "0.elements.0.elements.1.elements.0.elements"
	pathItemElement   = "elements.0"
	pathTrackName     = "attributes.aria-label"

	urlSearch      = "https://soundcloud.com/search?q=%s %s"
	urlRecommended = "https://soundcloud.com%s/recommended"

	trackTitleStart = `Track: `
	trackSeparator  = ` by `
	space           = ` `
	empty           = ""
)

type ctxt struct {
	ctx    context.Context
	cancel context.CancelFunc
}

type parser struct {
	parentCtx ctxt
	log       customLogger.Logger
}

func newParser(log customLogger.Logger) parser {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.Flag("disable-web-security", "1"),
	)

	allocatorCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocatorCtx)
	chromedp.Run(ctx)

	return parser{
		parentCtx: ctxt{
			ctx:    ctx,
			cancel: cancel,
		},
		log: log,
	}
}

func (p parser) getSimilar(artist, song string) []datastruct.AudioItem {
	ctxRelatedSongs, _ := chromedp.NewContext(p.parentCtx.ctx)
	ctxRelatedSongs, cancelRelatedSongs := context.WithTimeout(ctxRelatedSongs, 15*time.Second)
	defer cancelRelatedSongs()
	ctxUrl, cancelUrl := chromedp.NewContext(ctxRelatedSongs)
	defer cancelUrl()

	return p.getRelatedTracks(
		fmt.Sprintf(urlRecommended, p.getTrackPathname(artist, song, ctxt{ctxUrl, cancelUrl})),
		ctxt{ctxRelatedSongs, cancelRelatedSongs})
}

func (p parser) getTrackPathname(artist, song string, ctxt ctxt) string {
	var (
		nodes []*cdp.Node
		tasks = make(chromedp.Tasks, 4)
	)

	for i := 0; i < 4; i++ {
		c := i + 1
		tasks[i] = chromedp.ActionFunc(func(ctx context.Context) error {
			chromedp.Nodes(fmt.Sprintf(pathURLPathname, c), &nodes, chromedp.AtLeast(0)).Do(ctx)
			if len(nodes) == 0 {
				return nil
			}
			if !p.isASCII(artist) || !p.isASCII(song) {
				defer ctxt.cancel()
				return nil
			}
			if !p.match(song, nodes[0].Attributes[3]) {
				return nil
			}
			defer ctxt.cancel()
			return nil
		})
	}

	chromedp.Run(ctxt.ctx, chromedp.Navigate(fmt.Sprintf(urlSearch, artist, song)), tasks)

	if len(nodes) == 0 {
		return empty
	}

	return nodes[0].Attributes[3]
}

func (p parser) getRelatedTracks(trackRecommendURL string, ctxt ctxt) []datastruct.AudioItem {
	result := []datastruct.AudioItem{}

	if trackRecommendURL == empty {
		return result
	}

	var (
		nodes []*cdp.Node
		data  string
	)

	action := chromedp.ActionFunc(func(ctx context.Context) error {
		chromedp.Nodes(pathFlexContainer, &nodes, chromedp.AtLeast(0)).Do(ctx)

		if len(nodes) != 0 {
			defer time.Sleep(470 * time.Millisecond)
			return chromedp.Click(pathFlexContainer).Do(ctx)
		} else {
			defer ctxt.cancel()
			return chromedp.OuterHTML(pathItemList, &data, chromedp.AtLeast(0)).Do(ctx)
		}
	})

	tasks := make(chromedp.Tasks, 10)
	for i := 0; i < 10; i++ {
		tasks[i] = action
	}

	chromedp.Run(ctxt.ctx, chromedp.Navigate(trackRecommendURL), tasks)

	if data == empty {
		return result
	}

	gjson.GetBytes(p.reformat(data), pathItem).ForEach(func(key, value gjson.Result) bool {
		result = append(result, p.conversion(value.Get(pathItemElement).Get(pathTrackName).String()))
		return true
	})

	return result
}

func (p parser) CloseBrowser() {
	p.parentCtx.cancel()
}

func (p parser) reformat(htmlData string) []byte {
	d, err := html2json.New(strings.NewReader(htmlData))
	if err != nil {
		p.log.Warn(p.log.CallInfoStr(), err.Error())
	}

	js, err := d.ToJSON()
	if err != nil {
		p.log.Warn(p.log.CallInfoStr(), err.Error())
	}

	return js
}

func (p parser) conversion(trackStr string) datastruct.AudioItem {
	b := strings.Split(strings.Trim(trackStr, trackTitleStart), trackSeparator)
	if len(b) > 2 {
		b[0] = strings.Join(b[0:len(b)-2], space)
	}

	for _, separator := range separators {
		if strings.Contains(b[0], separator) {
			temp := strings.Split(b[0], separator)
			return datastruct.AudioItem{
				Artist: r.Replace(temp[0]),
				Title:  r.Replace(temp[1]),
			}
		}
	}

	return datastruct.AudioItem{
		Artist: r.Replace(b[len(b)-1]),
		Title:  r.Replace(b[0]),
	}
}

func (p parser) isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func (p parser) match(pattern, s string) bool {
	return strings.Contains(s, strings.ToLower(strings.ReplaceAll(pattern, ` `, `-`)))
}

var (
	r = strings.NewReplacer(
		"&#x27;", "`",
		"&quot;", `"`,
		"&amp;", "&",
	)
	separators = []string{
		` - `,
		` – `,
		` — `,
	}
)
