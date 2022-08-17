package gnoosic

import (
	"errors"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
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
	q := req.URL.Query()
	q.Add("skip", "1")
	q.Add("Fave01", emp)
	q.Add("Fave02", emp)
	q.Add("Fave03", emp)
	req.URL.RawQuery = q.Encode()
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
	url, _ := resp.Location()

	return strings.ReplaceAll(strings.TrimPrefix(url.Path, pArtist), `+`, ` `)
}
