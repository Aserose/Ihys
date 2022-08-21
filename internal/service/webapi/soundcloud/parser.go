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
	pURLPathname   = "#content > div > div > div.l-main > div > div > div > ul > li:nth-child(%d) > div > div > div > div.sound__artwork.sc-mr-1x > a"
	pFlexContainer = "#content > div > div > div.l-main.g-main-scroll-area > div > div > div.relatedList__list > div"
	pItemList      = "#content > div > div > div.l-main.g-main-scroll-area > div > div > div.relatedList__list > ul"
	pItem          = "0.elements.0.elements.1.elements.0.elements"
	pItemElement   = "elements.0"
	pTrackName     = "attributes.aria-label"

	urlSearch         = "https://soundcloud.com/search?q="
	urlRecommended    = "https://soundcloud.com%s/recommended"
	urlRecommendedEmp = "https://soundcloud.com/recommended"

	trackTitleStart = `Track: `
	trackSeparator  = ` by `
	spc             = ` `
	emp             = ""
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
		chromedp.Headless, chromedp.NoSandbox, chromedp.DisableGPU,
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.Flag("disable-web-security", "1"),
		chromedp.Flag("disable-setuid-sandbox", true),
	)

	alc, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(alc)
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

func (p parser) similar(artist, title string) []datastruct.Song {
	ctxRelatedSongs, cancelRS := chromedp.NewContext(p.parentCtx.ctx)
	defer cancelRS()
	ctxUrl, cancelUrl := chromedp.NewContext(ctxRelatedSongs)
	defer cancelUrl()

	return p.related(fmt.Sprintf(urlRecommended, p.trackPathname(artist, title, ctxt{ctxUrl, cancelUrl})), ctxt{ctxRelatedSongs, cancelRS})
}

func (p parser) trackPathname(artist, title string, ctxt ctxt) string {
	var (
		nodes []*cdp.Node
		tasks = make(chromedp.Tasks, 10)
	)

	childPathname := func(nthChildNum int) actionFunc {
		return func(ctx context.Context) error {
			if err := chromedp.Nodes(fmt.Sprintf(pURLPathname, nthChildNum), &nodes, chromedp.AtLeast(0)).Do(ctx); err != nil {
				p.log.Error(p.log.CallInfoStr(), err.Error())
			}

			if len(nodes) == 0 {
				return nil
			}
			if !p.isASCII(artist) || !p.isASCII(title) {
				defer ctxt.cancel()
				return nil
			}
			if !p.match(title, nodes[0].Attributes[3]) {
				return nil
			}

			defer ctxt.cancel()
			return nil
		}
	}

	for i := 0; i < 10; i++ {
		tasks[i] = chromedp.ActionFunc(childPathname(i))
	}

	wait := chromedp.ActionFunc(func(ctx context.Context) error {
		time.Sleep(276 * time.Millisecond)
		return nil
	})

	if err := chromedp.Run(ctxt.ctx, chromedp.Navigate(urlSearch+url.QueryEscape(artist+` `+title)), wait, tasks); err != nil {
		p.log.Error(p.log.CallInfoStr(), err.Error())
	}

	if len(nodes) == 0 {
		return emp
	}

	return nodes[0].Attributes[3]
}

func (p parser) related(rcmdURL string, ctxt ctxt) []datastruct.Song {
	res := []datastruct.Song{}

	if rcmdURL == urlRecommendedEmp {
		return res
	}

	nodes := []*cdp.Node{}
	var data string

	action := chromedp.ActionFunc(func(ctx context.Context) error {
		chromedp.Nodes(pFlexContainer, &nodes, chromedp.AtLeast(0)).Do(ctx)

		if len(nodes) != 0 {
			defer time.Sleep(470 * time.Millisecond)
			return chromedp.Click(pFlexContainer).Do(ctx)
		} else {
			defer ctxt.cancel()
			return chromedp.OuterHTML(pItemList, &data, chromedp.AtLeast(0)).Do(ctx)
		}
	})

	tasks := make(chromedp.Tasks, 10)
	for i := 0; i < 10; i++ {
		tasks[i] = action
	}

	wait := chromedp.ActionFunc(func(ctx context.Context) error {
		time.Sleep(200 * time.Millisecond)
		return nil
	})

	chromedp.Run(ctxt.ctx, chromedp.Navigate(rcmdURL), wait, tasks)

	if data == emp {
		return res
	}

	gjson.GetBytes(p.json(data), pItem).ForEach(func(key, value gjson.Result) bool {
		res = append(res, p.convert(value.Get(pItemElement).Get(pTrackName).String()))
		return true
	})

	return res
}

func (p parser) CloseBrowser() {
	p.parentCtx.cancel()
}

func (p parser) json(html string) []byte {
	d, err := html2json.New(strings.NewReader(html))
	if err != nil {
		p.log.Warn(p.log.CallInfoStr(), err.Error())
	}

	js, err := d.ToJSON()
	if err != nil {
		p.log.Warn(p.log.CallInfoStr(), err.Error())
	}

	return js
}

func (p parser) convert(track string) datastruct.Song {
	b := strings.Split(strings.Trim(track, trackTitleStart), trackSeparator)
	if len(b) > 2 {
		b[0] = strings.Join(b[0:len(b)-2], spc)
	}

	for _, sp := range separators {
		if strings.Contains(b[0], sp) {
			temp := strings.Split(b[0], sp)
			return datastruct.Song{
				Artist: replHtmlEnt.Replace(temp[0]),
				Title:  replHtmlEnt.Replace(temp[1]),
			}
		}
	}

	return datastruct.Song{
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
