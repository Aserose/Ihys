package genius

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
	"github.com/goccy/go-json"
	"github.com/valyala/fasthttp"
	"strings"
)

const (
	pSearch = `/search`

	urlBase   = `https://api.genius.com`
	urlSearch = urlBase + pSearch

	q = `q`

	hAuth = `AUTHORIZATION`
	vAuth = `Bearer `

	emp = ``
)

type enq struct {
	client *fasthttp.Client
	cfg    config.Genius
	log    customLogger.Logger
}

func newEnq(log customLogger.Logger, cfg config.Genius) enq {
	return enq{
		client: &fasthttp.Client{},
		cfg:    cfg,
		log:    log,
	}
}

func (e enq) send(req *fasthttp.Request) []byte {
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := e.client.Do(req, resp)
	if err != nil {
		e.log.Warn(e.log.CallInfoStr(), err.Error())
	}

	return resp.Body()
}

func (e enq) lyricsURL(audio datastruct.Song) string {
	resp := datastruct.GeniusSearch{}

	uri := fasthttp.AcquireURI()
	uri.Parse(nil, []byte(urlSearch))
	uri.QueryArgs().Add(q, audio.FirstArtist()+` `+audio.Title)

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set(hAuth, vAuth+e.cfg.Key)

	req.SetURI(uri)
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseURI(uri)
	}()
	json.Unmarshal(e.send(req), &resp)

	if len(resp.Response.Hits) == 0 {
		return emp
	}

	title := strings.ToLower(audio.Title)
	for _, hit := range resp.Response.Hits {
		if strings.Contains(strings.ToLower(hit.Result.Title), title) {
			if hit.Result.LyricsState != `complete` {
				return emp
			}
			return hit.Result.URL
		}
	}

	return emp
}
