package datastruct

type YaSongPage struct {
	Elements []struct {
		Elements []struct {
			Elements []struct {
				Text string `json:"text"`
			} `json:"elements"`
		} `json:"elements"`
	} `json:"elements"`
}

type YaSimilar struct {
	YaSidebar `json:"sidebarData"`
}

type YaSidebar struct {
	SimilarTracks []YaSimilarTracks `json:"similarTracks"`
}

type YaSimilarTracks struct {
	Artists []YaArtists `json:"artists"`
	Title   string      `json:"title"`
}

type YaArtists struct {
	Name string `json:"name"`
}
