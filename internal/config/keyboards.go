package config

type Buttons struct {
	Menus      `yaml:"menus"`
	LastFmBtn  Button `yaml:"last_fm"`
	YaMusicBtn Button `yaml:"ya_music"`
}

type Menus struct {
	MainMenu   Button `yaml:"main_menu"`
	SearchMenu Button `yaml:"search_menu"`
	SongMenu   `yaml:"song_menu"`
	VkSubMenu  `yaml:"vk_submenu"`
}

type SongMenu struct {
	Delete    Button `yaml:"delete"`
	Similiars Button `yaml:"similiars"`
	Best      Button `yaml:"best"`
}

type VkSubMenu struct {
	Self           Button `yaml:"vk_submenu_button"`
	Auth           Button `yaml:"vk_auth"`
	Recommendation Button `yaml:"vk_recommendation"`
	UserPlaylist   Button `yaml:"vk_user_playlist"`
}

type Button struct {
	Text         string `yaml:"text"`
	CallbackData string `yaml:"callback_data"`
	Delete       string `yaml:"delete"`
}
