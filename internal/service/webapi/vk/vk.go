package vk

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/repository"
	"IhysBestowal/internal/service/webapi/vk/client"
	"IhysBestowal/pkg/customLogger"
	"fmt"
	"github.com/goccy/go-json"
	"net/http"
)

const (
	getUserPlaylists         = "https://api.vk.com/method/audio.getPlaylists?access_token=%s&owner_id=%d&v=5.95"
	getPlaylistSong          = "https://api.vk.com/method/audio.get?access_token=%s&album_id=%d&owner_id=%d&v=5.95"
	getRecommendations       = "https://api.vk.com/method/audio.getRecommendations?access_token=%s&offset=%d&v=5.95"
	getRecommendationsCustom = "https://api.vk.com/method/audio.get?access_token=%s&album_id=-22&owner_id=%d&v=5.95"
	getAudioById             = "https://api.vk.com/method/audio.getById?access_token=%s&audios=%s&v=5.95"
	getUser                  = "https://api.vk.com/method/users.get?access_token=%s&v=5.95"
)

type httpClient interface {
	Send(req *http.Request) []byte
}

type VK struct {
	VAuth
	http httpClient
	log  customLogger.Logger
}

func New(log customLogger.Logger, cfg config.Vk, auth repository.Key) VK {
	c := client.New(log)
	return VK{
		VAuth: newVkAuth(log, cfg, auth, c),
		http:  c,
		log:   log,
	}
}

func (v VK) PlaylistSongs(user dto.TGUser, playlistId, ownerId int) (datastruct.Songs, error) {
	token, err := v.VAuth.token(user)
	if err != nil {
		return datastruct.Songs{}, err
	}
	result := datastruct.VKAudio{}

	getSongsID := fmt.Sprintf(getPlaylistSong, token, playlistId, ownerId)
	req, err := http.NewRequest(http.MethodGet, getSongsID, nil)
	if err != nil {
		v.log.Warn(v.log.CallInfoStr(), err.Error())
	}
	err = json.Unmarshal(v.Send(req), &result)
	if err != nil {
		v.log.Warn(v.log.CallInfoStr(), err.Error())
	}

	return v.newSongs(v.getById(v.songIds(result), token)), nil
}

func (v VK) UserPlaylists(user dto.TGUser) (datastruct.Playlists, error) {
	token, err := v.VAuth.token(user)
	if err != nil {
		return datastruct.Playlists{}, err
	}
	vkPlaylist := datastruct.VKPlaylist{}

	url := fmt.Sprintf(getUserPlaylists, token, v.VAuth.userId(token))
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		v.log.Warn(v.log.CallInfoStr(), err.Error())
	}
	err = json.Unmarshal(v.Send(req), &vkPlaylist)
	if err != nil {
		v.log.Warn(v.log.CallInfoStr(), err.Error())
	}

	return v.playlists(vkPlaylist), nil
}

func (v VK) playlists(pl datastruct.VKPlaylist) (playlist datastruct.Playlists) {
	playlist.From = "vk"
	playlist.Playlists = make([]datastruct.Playlist, len(pl.Response.Items))
	for i, p := range pl.Response.Items {
		playlist.Playlists[i].ID = p.ID
		playlist.Playlists[i].Title = p.Title
		playlist.Playlists[i].OwnerId = p.OwnerID
	}

	return
}

func (v VK) Recommendations(user dto.TGUser, offset int) (datastruct.Songs, error) {
	token, err := v.VAuth.token(user)
	if err != nil {
		return datastruct.Songs{}, err
	}
	result := datastruct.VKAudio{}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(getRecommendations, token, offset), nil)
	if err != nil {
		v.log.Warn(v.log.CallInfoStr(), err.Error())
	}

	err = json.Unmarshal(v.Send(req), &result)
	if err != nil {
		v.log.Warn(v.log.CallInfoStr(), err.Error())
	}

	return v.newSongs(result), nil
}

func (v VK) newSongs(res datastruct.VKAudio) (audio datastruct.Songs) {
	audio.Songs = make([]datastruct.Song, len(res.VKResponse.Items))

	for i, r := range res.VKResponse.Items {
		audio.Songs[i].Artist = r.Artist
		audio.Songs[i].Title = r.Title
		audio.Songs[i].Url = r.Url
	}

	return
}

func (v VK) RecommendationsCustom(user dto.TGUser) (datastruct.Songs, error) {
	token, err := v.VAuth.token(user)
	if err != nil {
		return datastruct.Songs{}, err
	}
	result := datastruct.VKAudio{}

	url := fmt.Sprintf(getRecommendationsCustom, token, v.VAuth.userId(token))

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		v.log.Warn(v.log.CallInfoStr(), err.Error())
	}

	err = json.Unmarshal(v.Send(req), &result)
	if err != nil {
		v.log.Warn(v.log.CallInfoStr(), err.Error())
	}

	return v.newSongs(v.getById(v.songIds(result), token)), nil
}

func (v VK) songIds(audio datastruct.VKAudio) (audioIds string) {
	for _, id := range audio.VKResponse.Items {
		audioIds += id.VKAds.ContentID + ","
	}
	return
}

func (v VK) getById(audioIds, token string) datastruct.VKAudio {
	result := datastruct.VKAudio{}

	getById := fmt.Sprintf(getAudioById, token, audioIds)
	req, err := http.NewRequest(http.MethodGet, getById, nil)
	if err != nil {
		v.log.Warn(v.log.CallInfoStr(), err.Error())
	}
	err = json.Unmarshal(v.Send(req), &result)
	if err != nil {
		v.log.Warn(v.log.CallInfoStr(), err.Error())
	}

	return result
}
