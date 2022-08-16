package auth

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/repository"
	"IhysBestowal/pkg/customLogger"
)

type Auth struct {
	vk Key
}

func New(log customLogger.Logger, cfg config.Auth, repo repository.Repository) Auth {
	return Auth{
		vk: newKey(log, cfg, repo.Vk()),
	}
}

func (a Auth) Vk() repository.Key { return a.vk }
