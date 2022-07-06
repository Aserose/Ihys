package datastruct

type GeniusSearch struct {
	Response struct {
		Hits []struct {
			Result struct {
				ArtistNames string `json:"artist_names"`
				LyricsState string `json:"lyrics_state"`
				Title       string `json:"title"`
				URL         string `json:"url"`
			} `json:"result"`
		} `json:"hits"`
	} `json:"response"`
}
