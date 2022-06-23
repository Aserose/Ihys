package discogs

type collater struct{}

func newCollater() collater {
	return collater{}
}

func (c collater) getFirstWebsite(urls []string) string {
	if len(urls) == 0 {
		return empty
	}
	return urls[0]
}
