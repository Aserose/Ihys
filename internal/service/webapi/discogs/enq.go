package discogs

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/pkg/customLogger"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/valyala/fasthttp"
	"strings"
)

const (
	urlBase   = "https://api.discogs.com"
	urlSearch = urlBase + pSearch

	pSearch = "/database/search"

	hAuth = `AUTHORIZATION`
	vAuth = "Discogs key=%s, secret=%s"

	q     = `q`
	qType = `type`

	typeArtist  = `artist`
	typeRelease = `release`
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
	uri.QueryArgs().Add(q, artist+` `+s.Title)
	uri.QueryArgs().Add(qType, typeRelease)

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set(hAuth, fmt.Sprintf(vAuth, e.cfg.Key, e.cfg.Secret))

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
	if query == emp {
		return []string{}
	}

	resource := e.URL(query, typeArg)
	if resource == emp {
		return nil
	}

	resp := datastruct.DiscogsResourceURL{}

	uri := fasthttp.AcquireURI()
	uri.Parse(nil, []byte(resource))

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set(hAuth, fmt.Sprintf(vAuth, e.cfg.Key, e.cfg.Secret))

	req.SetURI(uri)
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseURI(uri)
	}()
	json.Unmarshal(e.send(req), &resp)

	return resp.Websites
}

func (e enq) URL(query string, typeArg string) string {
	if query == emp {
		return query
	}

	resp := datastruct.DiscogsSearch{}

	uri := fasthttp.AcquireURI()
	uri.Parse(nil, []byte(urlSearch))
	if typeArg != emp {
		uri.QueryArgs().Add(qType, typeArg)
	}
	uri.QueryArgs().Add(q, query)

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set(hAuth, fmt.Sprintf(vAuth, e.cfg.Key, e.cfg.Secret))

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
	return emp
}
