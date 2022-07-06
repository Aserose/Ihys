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
	searchResp := datastruct.DiscogsSearch{}
	releaseResp := datastruct.DiscogsRelease{}
	artist := strings.ToLower(audio.GetFirstArtist())

	uri := fasthttp.AcquireURI()
	uri.Parse(nil, []byte(urlSearch))
	uri.QueryArgs().Add(queryQ, artist+` `+audio.Title)
	uri.QueryArgs().Add(queryType, typeRelease)

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set(headAuthKey, headAuthValueKey+e.cfg.Key+headAuthValueSecret+e.cfg.Secret)

	req.SetURI(uri)
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseURI(uri)
	}()
	json.Unmarshal(e.sendRequest(req), &searchResp)

	if len(searchResp.Results) == 0 {
		return datastruct.AudioInfo{}
	}

	for _, result := range searchResp.Results {

		if strings.Contains(strings.ToLower(result.Title), artist) {
			uri.Parse(nil, []byte(result.ResourceURL))
			req.SetURI(uri)
			json.Unmarshal(e.sendRequest(req), &releaseResp)

			return datastruct.AudioInfo{
				Label:       strings.Join(result.Label, ` | `),
				Genres:      append(releaseResp.Styles, releaseResp.Genres...),
				Country:     releaseResp.Country,
				ReleaseDate: releaseResp.ReleasedFormatted,
			}

		}

	}

	return datastruct.AudioInfo{}
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

	return resp.Websites
}

func (e enquirer) getResourceURL(query string, typeArg string) string {
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
	json.Unmarshal(e.sendRequest(req), &resp)

	query = strings.ToLower(query)
	for _, result := range resp.Results {
		if strings.Contains(strings.ToLower(result.Title), query) {
			return result.ResourceURL
		}
	}
	return empty
}
