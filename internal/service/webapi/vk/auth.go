package vk

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/service/auth"
	"IhysBestowal/pkg/customLogger"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"net/http"
)

type IAuth interface {
	Authorize(user dto.TGUser, serviceKey string) error
	GetAuthURL() string
	IsAuthorized(user dto.TGUser) bool
	TokenIsValid(token string) bool
	getKey(user dto.TGUser) (string, error)
	getUserId(token string) int
}

type vkAuth struct {
	auth        auth.IKey
	sendRequest func(req *http.Request) []byte
	authLink    string
	log         customLogger.Logger
}

func newVkAuth(log customLogger.Logger, cfg config.Vk, auth auth.IKey, sendRequest func(req *http.Request) []byte) IAuth {
	return vkAuth{
		auth:        auth,
		sendRequest: sendRequest,
		authLink:    cfg.AuthLink,
		log:         log,
	}
}

func (v vkAuth) getKey(user dto.TGUser) (string, error) {
	token := v.auth.GetKey(user)
	if !v.TokenIsValid(token) {
		v.auth.DeleteKey(user)
		return "", errors.New("invalid VK accessToken")
	}
	return token, nil
}

func (v vkAuth) TokenIsValid(token string) bool {
	var err = struct {
		Error struct {
			ErrorMsg string `json:"error_msg"`
		} `json:"error"`
	}{}

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(getUser, token), nil)
	json.Unmarshal(v.sendRequest(req), &err)

	return err.Error.ErrorMsg == ""
}

func (v vkAuth) GetAuthURL() string {
	return v.authLink
}

func (v vkAuth) Authorize(user dto.TGUser, serviceKey string) error {
	if !v.TokenIsValid(serviceKey) {
		return errors.New("invalid serviceKey")
	}

	v.auth.PutKey(user, serviceKey)

	return nil
}

func (v vkAuth) IsAuthorized(user dto.TGUser) bool {
	return v.TokenIsValid(v.auth.GetKey(user))
}

func (v vkAuth) getUserId(token string) int {
	resp := struct {
		Response []struct {
			Id int `json:"id"`
		} `json:"response"`
	}{}

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(getUser, token), nil)
	json.Unmarshal(v.sendRequest(req), &resp)

	return resp.Response[0].Id
}
