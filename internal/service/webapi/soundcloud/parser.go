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
	"net/url"
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

	urlSearch           = "https://soundcloud.com/search?q="
	urlRecommended      = "https://soundcloud.com%s/recommended"
	urlRecommendedEmpty = "https://soundcloud.com/recommended"

	trackTitleStart = `Track: `
	trackSeparator  = ` by `
	space           = ` `
	empty           = ""
)

type actionFunc func(ctx context.Context) error

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
		chromedp.Headless,
		chromedp.NoSandbox,
		chromedp.DisableGPU,
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.Flag("disable-web-security", "1"),
		chromedp.Flag("disable-setuid-sandbox", true),
	)

	allocatorCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocatorCtx)
	if err := chromedp.Run(ctx); err != nil {
		log.Error(log.CallInfoStr(), err.Error())
	}

	return parser{
		parentCtx: ctxt{
			ctx:    ctx,
			cancel: cancel,
		},
		log: log,
	}
}

func (p parser) getSimilar(artist, song string) []datastruct.AudioItem {
	ctxRelatedSongs, cancelRelatedSongs := chromedp.NewContext(p.parentCtx.ctx)
	defer cancelRelatedSongs()
	ctxUrl, cancelUrl := chromedp.NewContext(ctxRelatedSongs)
	defer cancelUrl()

	return p.getRelatedTracks(fmt.Sprintf(urlRecommended, p.getTrackPathname(artist, song, ctxt{ctxUrl, cancelUrl})), ctxt{ctxRelatedSongs, cancelRelatedSongs})
}

func (p parser) getTrackPathname(artist, song string, ctxt ctxt) string {
	var (
		nodes []*cdp.Node
		tasks = make(chromedp.Tasks, 10)
	)

	getPathnameFromChild := func(nthChildNum int) actionFunc {
		return func(ctx context.Context) error {
			if err := chromedp.Nodes(fmt.Sprintf(pathURLPathname, nthChildNum), &nodes, chromedp.AtLeast(0)).Do(ctx); err != nil {
				p.log.Error(p.log.CallInfoStr(), err.Error())
			}

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
		}
	}

	for i := 0; i < 10; i++ {
		tasks[i] = chromedp.ActionFunc(getPathnameFromChild(i))
	}

	wait := chromedp.ActionFunc(func(ctx context.Context) error {
		time.Sleep(266 * time.Millisecond)
		return nil
	})

	if err := chromedp.Run(ctxt.ctx, chromedp.Navigate(urlSearch+url.QueryEscape(artist+` `+song)), wait, tasks); err != nil {
		p.log.Error(p.log.CallInfoStr(), err.Error())
	}

	if len(nodes) == 0 {
		return empty
	}

	return nodes[0].Attributes[3]
}

func (p parser) getRelatedTracks(trackRecommendURL string, ctxt ctxt) []datastruct.AudioItem {
	result := []datastruct.AudioItem{}

	if trackRecommendURL == urlRecommendedEmpty {
		return result
	}

	nodes := []*cdp.Node{}
	var data string

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

	wait := chromedp.ActionFunc(func(ctx context.Context) error {
		time.Sleep(193 * time.Millisecond)
		return nil
	})

	chromedp.Run(ctxt.ctx, chromedp.Navigate(trackRecommendURL), wait, tasks)

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
				Artist: replHtmlEnt.Replace(temp[0]),
				Title:  replHtmlEnt.Replace(temp[1]),
			}
		}
	}

	return datastruct.AudioItem{
		Artist: replHtmlEnt.Replace(b[len(b)-1]),
		Title:  replHtmlEnt.Replace(b[0]),
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
	return strings.Contains(replSpecSym.Replace(strings.ToLower(s)), strings.ToLower(strings.ReplaceAll(pattern, ` `, `-`)))
}

var (
	replSpecSym = strings.NewReplacer(
		"'", ``,
	)
	replHtmlEnt = strings.NewReplacer(
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
