package auth

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/repository"
	"IhysBestowal/pkg/customLogger"
)

type Auth struct {
	repo repository.Repository
	cfg  config.Auth
	log  customLogger.Logger
}

func New(log customLogger.Logger, cfg config.Auth, repo repository.Repository) Auth {
	return Auth{
		repo: repo,
		cfg:  cfg,
		log:  log,
	}
}

func (as Auth) Vk() repository.Key {
	return newKey(as.log, as.cfg, as.repo.Vk())
}
