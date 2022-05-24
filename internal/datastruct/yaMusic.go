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
	Sidebar struct {
		SimilarTracks []struct {
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
			Title string `json:"title"`
		} `json:"similarTracks"`
	} `json:"sidebarData"`
}
