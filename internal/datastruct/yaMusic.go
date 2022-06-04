package datastruct

type YaMSourcePage struct {
	Elements []struct {
		Elements []struct {
			Elements []struct {
				Text string `json:"text"`
			} `json:"elements"`
		} `json:"elements"`
	} `json:"elements"`
}

type YaMSimiliar struct {
	YaMSidebar `json:"sidebarData"`
}

type YaMSidebar struct {
	SimilarTracks []YaMSimilarTracks `json:"similarTracks"`
}

type YaMSimilarTracks struct {
	Artists []YaMArtists `json:"artists"`
	Title   string       `json:"title"`
}

type YaMArtists struct {
	Name string `json:"name"`
}
