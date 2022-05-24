package datastruct

type LastFMAudio struct {
	Response []LastFMResponse `json:"response"`
}

type LastFMResponse struct {
	Artist string `json:"artist"`
	Title  string `json:"title"`
	Url    string `json:"url"`
}
