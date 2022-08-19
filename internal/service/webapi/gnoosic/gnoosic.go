package gnoosic

import (
	"errors"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

const (
	hHost  = `Host`
	hCache = `Cache-Control`
	vHost  = `www.gnoosic.com`
	vCache = `no-cache`

	pFaves  = "/faves.php"
	pArtist = "/artist/"

	urlBase   = "https://www.gnoosic.com"
	urlFront  = urlBase + pFaves
	urlArtist = urlBase + pArtist

	errRedirect = `redirect`
	emp         = ``
)

type Gnoosic struct {
	client *http.Client
}

func New() Gnoosic {
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
	req.URL.RawQuery = url.Values{
		"skip":   {"1"},
		"Fave01": {emp},
		"Fave02": {emp},
		"Fave03": {emp},
	}.Encode()

	req.Header.Set(hCache, vCache)
	req.Header.Set(hHost, vHost)

	client.Do(req)

	return Gnoosic{
		client: client,
	}
}

func (g Gnoosic) RandomArtist() string {
	req, _ := http.NewRequest(http.MethodGet, urlArtist, nil)
	req.Header.Set(hHost, vHost)

	resp, _ := g.client.Do(req)
	lct, _ := resp.Location()

	return strings.ReplaceAll(strings.TrimPrefix(lct.Path, pArtist), `+`, ` `)
}
