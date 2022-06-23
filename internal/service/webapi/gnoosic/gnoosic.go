package gnoosic

import (
	"errors"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

const (
	headerHost  = `Host`
	headerCache = `Cache-Control`
	hostValue   = `www.gnoosic.com`
	cacheValue  = `no-cache`

	pathFaves  = "/faves.php"
	pathArtist = "/artist/"

	urlBase   = "https://www.gnoosic.com"
	urlFront  = urlBase + pathFaves
	urlArtist = urlBase + pathArtist

	errRedirect = `redirect`
	empty       = ``
)

type IGnoosic interface {
	GetRandomArtist() string
}

type gnoosic struct {
	client *http.Client
}

func NewGnoosic() IGnoosic {
	jar, _ := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})

	client := &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New(errRedirect)
		},
	}

	req, _ := http.NewRequest(http.MethodPost, urlFront, nil)
	q := req.URL.Query()
	q.Add("skip", "1")
	q.Add("Fave01", empty)
	q.Add("Fave02", empty)
	q.Add("Fave03", empty)
	req.URL.RawQuery = q.Encode()
	req.Header.Set(headerCache, cacheValue)
	req.Header.Set(headerHost, hostValue)

	client.Do(req)

	return gnoosic{
		client: client,
	}
}

func (g gnoosic) GetRandomArtist() string {
	req, _ := http.NewRequest(http.MethodGet, urlArtist, nil)
	req.Header.Set(headerHost, hostValue)

	resp, _ := g.client.Do(req)
	url, _ := resp.Location()

	return strings.ReplaceAll(strings.TrimPrefix(url.Path, pathArtist), `+`, ` `)
}
