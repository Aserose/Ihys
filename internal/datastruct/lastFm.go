package datastruct

type LFMUnmr struct {
	LFMSimilarTracks  `json:"similartracks"`
	LFMSimilarArtists `json:"similarartists"`
	LFMTopTracks      `json:"toptracks"`
}

type LFMSimilarArtists struct {
	Artists []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"artist"`
}

type LFMSimilarTracks struct {
	Tracks []struct {
		Artist struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"artist"`
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"track"`
}

type LFMTopTracks struct {
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

type LFMSearchTrack struct {
	Results struct {
		TrackMatches struct {
			Tracks []struct {
				Name   string `json:"name"`
				Artist string `json:"artist"`
				Url    string `json:"url"`
			} `json:"track"`
		} `json:"trackmatches"`
	} `json:"results"`
}
