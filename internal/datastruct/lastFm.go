package datastruct

type LastFMAudio struct {
	Response []LastFMResponse `json:"response"`
}

type LastFMResponse struct {
	Artist string `json:"artist"`
	Title  string `json:"title"`
	Url    string `json:"url"`
}

type LastFMUnmr struct {
	LastFMSimilarTracks  `json:"similartracks"`
	LastFMSimilarArtists `json:"similarartists"`
	LastFMTopTracks      `json:"toptracks"`
}

type LastFMSimilarArtists struct {
	Artists []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"artist"`
}

type LastFMSimilarTracks struct {
	Tracks []struct {
		Artist struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"artist"`
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"track"`
}

type LastFMTopTracks struct {
	Tracks []struct {
		Artist struct {
			Mbid string `json:"mbid"`
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"artist"`
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"track"`
}
