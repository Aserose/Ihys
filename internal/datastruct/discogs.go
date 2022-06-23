package datastruct

type DiscogsSearch struct {
	Results []struct {
		Title       string `json:"title"`
		ResourceURL string `json:"resource_url"`
	} `json:"results"`
}

type DiscogsResourceURL struct {
	URLs []string `json:"urls"`
}

type DiscogsRelease struct {
	Results []struct {
		Label   []string `json:"label"`
		Style   []string `json:"style"`
		Genre   []string `json:"genre"`
		Country string   `json:"country"`
		Year    string   `json:"year"`
	} `json:"results"`
}
