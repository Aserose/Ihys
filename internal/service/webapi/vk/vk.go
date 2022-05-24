package vk

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/service/auth"
	"IhysBestowal/pkg/customLogger"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type IVk interface {
	Auth() IAuth
	GetRecommendations(user dto.TGUser, offset int) (datastruct.AudioItems, error)
	GetRecommendationsCustom(user dto.TGUser) (datastruct.AudioItems, error)
	GetUserPlaylists(user dto.TGUser) (datastruct.PlaylistItems, error)
	GetPlaylistSongs(user dto.TGUser, playlistId, ownerId int) (datastruct.AudioItems, error)
}

type vk struct {
	auth        IAuth
	sendRequest func(req *http.Request) []byte
	log         customLogger.Logger
}

func NewVK(log customLogger.Logger, cfg config.Vk, auth auth.IKey) IVk {
	httpClient := &http.Client{}

	sendRequest := func(req *http.Request) []byte {
		req.Header.Set(
			`User-Agent`,
			`VKAndroidApp/4.13.1-1206 (Android 4.4.3; SDK 19; armeabi; ; ru)","Accept": "image/gif, image/x-xbitmap, image/jpeg, image/pjpeg, */*`)

		b, err := httpClient.Do(req)
		if err != nil {
			log.Warn(log.CallInfoStr(), err.Error())
		}
		body, err := io.ReadAll(b.Body)
		if err != nil {
			log.Warn(log.CallInfoStr(), err.Error())
		}

		return body
	}

	return vk{
		auth:        newVkAuth(log, cfg, auth, sendRequest),
		sendRequest: sendRequest,
		log:         log,
	}
}

func (v vk) Auth() IAuth {
	return v.auth
}

func (v vk) GetPlaylistSongs(user dto.TGUser, playlistId, ownerId int) (datastruct.AudioItems, error) {
	token, err := v.auth.getKey(user)
	if err != nil {
		return datastruct.AudioItems{}, err
	}
	result := datastruct.VKAudio{}

	getSongsID := fmt.Sprintf(getPlaylistSong, token, playlistId, ownerId)
	req, err := http.NewRequest(http.MethodGet, getSongsID, nil)
	if err != nil {
		v.log.Warn(v.log.CallInfoStr(), err.Error())
	}
	err = json.Unmarshal(v.sendRequest(req), &result)
	if err != nil {
		v.log.Warn(v.log.CallInfoStr(), err.Error())
	}

	return v.newAudioItems(v.getById(v.getAudioIds(result), token)), nil
}

func (v vk) GetUserPlaylists(user dto.TGUser) (datastruct.PlaylistItems, error) {
	token, err := v.auth.getKey(user)
	if err != nil {
		return datastruct.PlaylistItems{}, err
	}
	vkPlaylist := datastruct.VKPlaylist{}

	url := fmt.Sprintf(getUserPlaylists, token, v.auth.getUserId(token))
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		v.log.Warn(v.log.CallInfoStr(), err.Error())
	}
	err = json.Unmarshal(v.sendRequest(req), &vkPlaylist)
	if err != nil {
		v.log.Warn(v.log.CallInfoStr(), err.Error())
	}

	return v.newPlaylistItems(vkPlaylist), nil
}

func (v vk) newPlaylistItems(pl datastruct.VKPlaylist) (playlist datastruct.PlaylistItems) {
	playlist.From = "vk"
	playlist.Items = make([]datastruct.PlaylistItem, len(pl.Response.Items))
	for i, p := range pl.Response.Items {
		playlist.Items[i].ID = p.ID
		playlist.Items[i].Title = p.Title
		playlist.Items[i].OwnerId = p.OwnerID
	}

	return
}

func (v vk) GetRecommendations(user dto.TGUser, offset int) (datastruct.AudioItems, error) {
	token, err := v.auth.getKey(user)
	if err != nil {
		return datastruct.AudioItems{}, err
	}
	result := datastruct.VKAudio{}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(getRecommendantions, token, offset), nil)
	if err != nil {
		v.log.Warn(v.log.CallInfoStr(), err.Error())
	}

	err = json.Unmarshal(v.sendRequest(req), &result)
	if err != nil {
		v.log.Warn(v.log.CallInfoStr(), err.Error())
	}

	return v.newAudioItems(result), nil
}

func (v vk) newAudioItems(res datastruct.VKAudio) (audio datastruct.AudioItems) {
	audio.Items = make([]datastruct.AudioItem, len(res.VKResponse.Items))

	for i, r := range res.VKResponse.Items {
		audio.Items[i].Artist = r.Artist
		audio.Items[i].Title = r.Title
		audio.Items[i].Url = r.Url
	}

	return
}

func (v vk) GetRecommendationsCustom(user dto.TGUser) (datastruct.AudioItems, error) {
	token, err := v.auth.getKey(user)
	if err != nil {
		return datastruct.AudioItems{}, err
	}
	result := datastruct.VKAudio{}

	url := fmt.Sprintf(getRecommendantionsCustom, token, v.auth.getUserId(token))

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		v.log.Warn(v.log.CallInfoStr(), err.Error())
	}

	err = json.Unmarshal(v.sendRequest(req), &result)
	if err != nil {
		v.log.Warn(v.log.CallInfoStr(), err.Error())
	}

	return v.newAudioItems(v.getById(v.getAudioIds(result), token)), nil
}

func (v vk) getAudioIds(audio datastruct.VKAudio) (audioIds string) {
	for _, id := range audio.VKResponse.Items {
		audioIds += id.VKAds.ContentID + ","
	}
	return
}

func (v vk) getById(audioIds, token string) datastruct.VKAudio {
	result := datastruct.VKAudio{}

	getById := fmt.Sprintf(getAudioById, token, audioIds)
	req, err := http.NewRequest(http.MethodGet, getById, nil)
	if err != nil {
		v.log.Warn(v.log.CallInfoStr(), err.Error())
	}
	err = json.Unmarshal(v.sendRequest(req), &result)
	if err != nil {
		v.log.Warn(v.log.CallInfoStr(), err.Error())
	}

	return result
}
