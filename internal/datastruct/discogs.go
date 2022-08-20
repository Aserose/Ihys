package datastruct

type DiscogsSearch struct {
	Results []struct {
		Label       []string `json:"label"`
		Title       string   `json:"title"`
		ResourceURL string   `json:"resource_url"`
	} `json:"results"`
}

type DiscogsURL struct {
	Urls []string `json:"urls"`
}

type DiscogsRelease struct {
	Country string `json:"country"`
	Labels  []struct {
		Name string `json:"name"`
	} `json:"labels"`
	Released          string   `json:"released"`
	ReleasedFormatted string   `json:"released_formatted"`
	Genres            []string `json:"genres"`
	Styles            []string `json:"styles"`
}
