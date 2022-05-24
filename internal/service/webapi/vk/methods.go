package vk

var (
	getUserPlaylists          = "https://api.vk.com/method/audio.getPlaylists?access_token=%s&owner_id=%d&v=5.95"
	getPlaylistSong           = "https://api.vk.com/method/audio.get?access_token=%s&album_id=%d&owner_id=%d&v=5.95"
	getRecommendantions       = "https://api.vk.com/method/audio.getRecommendations?access_token=%s&offset=%d&v=5.95"
	getRecommendantionsCustom = "https://api.vk.com/method/audio.get?access_token=%s&album_id=-22&owner_id=%d&v=5.95"
	getAudioById              = "https://api.vk.com/method/audio.getById?access_token=%s&audios=%s&v=5.95"
	getUser                   = "https://api.vk.com/method/users.get?access_token=%s&v=5.95"
)
