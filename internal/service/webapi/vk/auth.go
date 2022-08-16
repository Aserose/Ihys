package vk

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/repository"
	"IhysBestowal/pkg/customLogger"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"net/http"
)

type VAuth struct {
	auth repository.Key
	httpClient
	authLink string
	log      customLogger.Logger
}

func newVkAuth(log customLogger.Logger, cfg config.Vk, auth repository.Key, client httpClient) VAuth {
	return VAuth{
		auth:       auth,
		httpClient: client,
		authLink:   cfg.AuthLink,
		log:        log,
	}
}

func (v VAuth) token(user dto.TGUser) (string, error) {
	token := v.auth.Get(user)
	if !v.IsValid(token) {
		v.auth.Delete(user)
		return "", errors.New("invalid VK token")
	}
	return token, nil
}

func (v VAuth) IsValid(token string) bool {
	var err = struct {
		Error struct {
			ErrorMsg string `json:"error_msg"`
		} `json:"error"`
	}{}

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(getUser, token), nil)
	json.Unmarshal(v.Send(req), &err)

	return err.Error.ErrorMsg == ""
}

func (v VAuth) AuthURL() string {
	return v.authLink
}

func (v VAuth) Auth(user dto.TGUser, serviceKey string) error {
	if !v.IsValid(serviceKey) {
		return errors.New("invalid serviceKey")
	}

	v.auth.Create(user, serviceKey)

	return nil
}

func (v VAuth) IsAuthorized(user dto.TGUser) bool {
	return v.IsValid(v.auth.Get(user))
}

func (v VAuth) userId(token string) int {
	resp := struct {
		Response []struct {
			Id int `json:"id"`
		} `json:"response"`
	}{}

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(getUser, token), nil)
	json.Unmarshal(v.Send(req), &resp)

	return resp.Response[0].Id
}
