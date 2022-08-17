package discogs

type clt struct{}

func newClt() clt {
	return clt{}
}

func (c clt) first(urls []string) string {
	if len(urls) == 0 {
		return emp
	}
	return urls[0]
}
