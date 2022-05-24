package datastruct

type VKAudio struct {
	VKResponse struct {
		Count int `json:"count,omitempty"`
		Items []struct {
			Artist string `json:"artist"`
			Title  string `json:"title"`
			Url    string `json:"url"`
			VKAds  struct {
				ContentID string `json:"content_id"`
			} `json:"ads"`
		} `json:"items"`
	} `json:"response"`
}

type VKPlaylist struct {
	Response struct {
		Count  int           `json:"count"`
		Groups []interface{} `json:"groups"`
		Items  []struct {
			AccessKey   string `json:"access_key"`
			AlbumType   string `json:"album_type"`
			Count       int    `json:"count"`
			CreateTime  int    `json:"create_time"`
			Description string `json:"description"`
			Followed    *struct {
				OwnerID    int `json:"owner_id"`
				PlaylistID int `json:"playlist_id"`
			} `json:"followed,omitempty"`
			Followers int `json:"followers"`
			Genres    []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"genres"`
			ID          int  `json:"id"`
			IsExplicit  bool `json:"is_explicit"`
			IsFollowing bool `json:"is_following"`
			MainArtists []struct {
				Domain string `json:"domain"`
				ID     string `json:"id"`
				Name   string `json:"name"`
			} `json:"main_artists,omitempty"`
			Original *struct {
				AccessKey  string `json:"access_key"`
				OwnerID    int    `json:"owner_id"`
				PlaylistID int    `json:"playlist_id"`
			} `json:"original,omitempty"`
			OwnerID int `json:"owner_id"`
			Photo   *struct {
				Height    int    `json:"height"`
				Photo1200 string `json:"photo_1200"`
				Photo135  string `json:"photo_135"`
				Photo270  string `json:"photo_270"`
				Photo300  string `json:"photo_300"`
				Photo34   string `json:"photo_34"`
				Photo600  string `json:"photo_600"`
				Photo68   string `json:"photo_68"`
				Width     int    `json:"width"`
			} `json:"photo,omitempty"`
			Plays    int    `json:"plays"`
			Subtitle string `json:"subtitle,omitempty"`
			Thumbs   []struct {
				Height    int    `json:"height"`
				Photo1200 string `json:"photo_1200"`
				Photo135  string `json:"photo_135"`
				Photo270  string `json:"photo_270"`
				Photo300  string `json:"photo_300"`
				Photo34   string `json:"photo_34"`
				Photo600  string `json:"photo_600"`
				Photo68   string `json:"photo_68"`
				Width     int    `json:"width"`
			} `json:"thumbs,omitempty"`
			Title      string `json:"title"`
			Type       int    `json:"type"`
			UpdateTime int    `json:"update_time"`
			Year       int    `json:"year,omitempty"`
		} `json:"items"`
		NextFrom string        `json:"next_from"`
		Profiles []interface{} `json:"profiles"`
	} `json:"response"`
}
