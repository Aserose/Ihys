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
	pathSearch = "/database/search"

	urlBase   = "https://api.discogs.com"
	urlSearch = urlBase + pathSearch

	typeArtist  = `artist`
	typeRelease = `release`

	discogsName         = `Discogs`
	setKey              = `key=`
	setSecret           = `secret=`
	headAuthKey         = `AUTHORIZATION`
	headAuthValueKey    = discogsName + ` ` + setKey
	headAuthValueSecret = `, ` + setSecret

	queryQ    = `q`
	queryType = `type`
)

type enq struct {
	httpClient *fasthttp.Client
	cfg        config.Discogs
	log        customLogger.Logger
}

func newEnq(log customLogger.Logger, cfg config.Discogs) enq {
	return enq{
		httpClient: &fasthttp.Client{},
		cfg:        cfg,
		log:        log,
	}
}

func (e enq) send(req *fasthttp.Request) []byte {
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := e.httpClient.Do(req, resp)
	if err != nil {
		e.log.Warn(e.log.CallInfoStr(), err.Error())
	}

	return resp.Body()
}

func (e enq) songInfo(s datastruct.Song) datastruct.SongInfo {
	searchResp := datastruct.DiscogsSearch{}
	releaseResp := datastruct.DiscogsRelease{}
	artist := strings.ToLower(s.FirstArtist())

	uri := fasthttp.AcquireURI()
	uri.Parse(nil, []byte(urlSearch))
	uri.QueryArgs().Add(queryQ, artist+` `+s.Title)
	uri.QueryArgs().Add(queryType, typeRelease)

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set(headAuthKey, headAuthValueKey+e.cfg.Key+headAuthValueSecret+e.cfg.Secret)

	req.SetURI(uri)
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseURI(uri)
	}()
	json.Unmarshal(e.send(req), &searchResp)

	if len(searchResp.Results) == 0 {
		return datastruct.SongInfo{}
	}

	for _, res := range searchResp.Results {

		if strings.Contains(strings.ToLower(res.Title), artist) {
			uri.Parse(nil, []byte(res.ResourceURL))
			req.SetURI(uri)
			json.Unmarshal(e.send(req), &releaseResp)

			return datastruct.SongInfo{
				Label:       strings.Join(res.Label, ` | `),
				Genres:      append(releaseResp.Styles, releaseResp.Genres...),
				Country:     releaseResp.Country,
				ReleaseDate: releaseResp.ReleasedFormatted,
			}

		}

	}

	return datastruct.SongInfo{}
}

func (e enq) sites(query string, typeArg string) []string {
	if query == empty {
		return []string{}
	}

	resource := e.URL(query, typeArg)
	if resource == empty {
		return nil
	}

	resp := datastruct.DiscogsResourceURL{}

	uri := fasthttp.AcquireURI()
	uri.Parse(nil, []byte(resource))

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set(headAuthKey, headAuthValueKey+e.cfg.Key+headAuthValueSecret+e.cfg.Secret)

	req.SetURI(uri)
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseURI(uri)
	}()
	json.Unmarshal(e.send(req), &resp)

	return resp.Websites
}

func (e enq) URL(query string, typeArg string) string {
	if query == empty {
		return query
	}

	resp := datastruct.DiscogsSearch{}

	uri := fasthttp.AcquireURI()
	uri.Parse(nil, []byte(urlSearch))
	if typeArg != empty {
		uri.QueryArgs().Add(queryType, typeArg)
	}
	uri.QueryArgs().Add(queryQ, query)

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set(headAuthKey, headAuthValueKey+e.cfg.Key+headAuthValueSecret+e.cfg.Secret)

	req.SetURI(uri)
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseURI(uri)
	}()
	json.Unmarshal(e.send(req), &resp)

	query = strings.ToLower(query)
	for _, result := range resp.Results {
		if strings.Contains(strings.ToLower(result.Title), query) {
			return result.ResourceURL
		}
	}
	return empty
}
