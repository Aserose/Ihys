package discogs

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
	"github.com/goccy/go-json"
	"github.com/valyala/fasthttp"
	"strings"
)

const (
	pathSearch = "database/search"

	urlBase   = "https://api.discogs.com/"
	urlSearch = urlBase + pathSearch

	typeArtist  = `artist`
	typeRelease = `release`

	discogsName         = `Discogs`
	setKey              = `key=`
	setSecret           = `secret=`
	headAuthKey         = `AUTHORIZATION`
	headAuthValueKey    = discogsName + ` ` + setKey
	headAuthValueSecret = `, ` + setSecret
)

type enquirer struct {
	httpClient *fasthttp.Client
	cfg        config.Discogs
	log        customLogger.Logger
}

func newEnquirer(log customLogger.Logger, cfg config.Discogs) enquirer {
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

func (e enquirer) getSongInfo(audio datastruct.AudioItem) datastruct.AudioInfo {
	resp := datastruct.DiscogsRelease{}

	uri := fasthttp.AcquireURI()
	uri.Parse(nil, []byte(urlSearch))
	uri.QueryArgs().Add(`q`, audio.GetFirstArtist()+` `+audio.Title)
	uri.QueryArgs().Add(`type`, typeRelease)

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set(headAuthKey, headAuthValueKey+e.cfg.Key+headAuthValueSecret+e.cfg.Secret)

	req.SetURI(uri)
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseURI(uri)
	}()
	json.Unmarshal(e.sendRequest(req), &resp)

	if len(resp.Results) == 0 {
		return datastruct.AudioInfo{}
	}

	return datastruct.AudioInfo{
		Label:   resp.Results[0].Label[0],
		Genres:  append(resp.Results[0].Style, resp.Results[0].Genre...),
		Country: resp.Results[0].Country,
		Year:    resp.Results[0].Year,
	}
}

func (e enquirer) getWebsites(query string, typeArg string) []string {
	if query == empty {
		return []string{}
	}

	resourceURL := e.getResourceURL(query, typeArg)
	if resourceURL == empty {
		return nil
	}

	resp := datastruct.DiscogsResourceURL{}

	uri := fasthttp.AcquireURI()
	uri.Parse(nil, []byte(resourceURL))

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set(headAuthKey, headAuthValueKey+e.cfg.Key+headAuthValueSecret+e.cfg.Secret)

	req.SetURI(uri)
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseURI(uri)
	}()
	json.Unmarshal(e.sendRequest(req), &resp)

	return resp.URLs
}

func (e enquirer) getResourceURL(query string, typeArg string) string {
	if query == empty {
		return query
	}

	resp := datastruct.DiscogsSearch{}

	uri := fasthttp.AcquireURI()
	uri.Parse(nil, []byte(urlSearch))
	if typeArg != empty {
		uri.QueryArgs().Add(`type`, typeArg)
	}
	uri.QueryArgs().Add(`q`, query)

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set(headAuthKey, headAuthValueKey+e.cfg.Key+headAuthValueSecret+e.cfg.Secret)

	req.SetURI(uri)
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseURI(uri)
	}()
	json.Unmarshal(e.sendRequest(req), &resp)

	query = strings.ToLower(query)
	for _, result := range resp.Results {
		if strings.Contains(strings.ToLower(result.Title), query) {
			return result.ResourceURL
		}
	}
	return empty
}
