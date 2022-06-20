package config

type Keypads struct {
	MainMenu   `yaml:"main_menu"`
	SearchMenu `yaml:"search_menu"`
	SongMenu   `yaml:"song_menu"`
}

type MainMenu struct {
	VkSubMenu `yaml:"vk_submenu"`
}

type SearchMenu struct {
	Self    Button `yaml:"search_menu_button"`
	YaMusic Button `yaml:"ya_music"`
	LastFM  Button `yaml:"last_fm"`
	All     Button `yaml:"all"`
}

type SongMenu struct {
	Self    Button `yaml:"song_menu_button"`
	Delete  Button `yaml:"delete"`
	Similar Button `yaml:"similar"`
	Best    Button `yaml:"best"`
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
}
