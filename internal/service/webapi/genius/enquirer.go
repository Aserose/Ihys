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
	pathSearch = `/search`

	urlBase   = `https://api.genius.com`
	urlSearch = urlBase + pathSearch

	query = `q`

	headAuthKey   = `AUTHORIZATION`
	headAuthValue = `Bearer `

	empty = ``
)

type enquirer struct {
	httpClient *fasthttp.Client
	cfg        config.Genius
	log        customLogger.Logger
}

func newEnquirer(log customLogger.Logger, cfg config.Genius) enquirer {
	return enquirer{
		httpClient: &fasthttp.Client{},
		cfg:        cfg,
		log:        log,
	}
}

func (e enquirer) sendRequest(req *fasthttp.Request) []byte {
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := e.httpClient.Do(req, resp)
	if err != nil {
		e.log.Warn(e.log.CallInfoStr(), err.Error())
	}

	return resp.Body()
}

func (e enquirer) getLyricsURL(audio datastruct.AudioItem) string {
	resp := datastruct.GeniusSearch{}

	uri := fasthttp.AcquireURI()
	uri.Parse(nil, []byte(urlSearch))
	uri.QueryArgs().Add(query, audio.GetFirstArtist()+` `+audio.Title)

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set(headAuthKey, headAuthValue+e.cfg.Key)

	req.SetURI(uri)
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseURI(uri)
	}()
	json.Unmarshal(e.sendRequest(req), &resp)

	if len(resp.Response.Hits) == 0 {
		return empty
	}

	songTitle := strings.ToLower(audio.Title)
	for _, hit := range resp.Response.Hits {
		if strings.Contains(strings.ToLower(hit.Result.Title), songTitle) {
			if hit.Result.LyricsState != `complete` {
				return empty
			}
			return hit.Result.URL
		}
	}

	return empty
}
