package service

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/repository"
	"IhysBestowal/internal/service/auth"
	"IhysBestowal/internal/service/menu"
	"IhysBestowal/internal/service/webapi"
	"IhysBestowal/pkg/customLogger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramService struct {
	menu.TGMenu
	SendMsg func(chattable tgbotapi.Chattable) tgbotapi.Message
}

type AuthService struct {
	auth.AuthService
	IsValidToken func(token string) bool //TODO
}

type Service struct {
	AuthService
	TelegramService
}

func NewService(log customLogger.Logger, cfg config.Service, repo repository.Repository) Service {
	authService := auth.NewAuthService(log, cfg.Auth, repo)
	webApiService := webapi.NewWebApiService(log, cfg, repo, authService)

	return Service{
		AuthService: AuthService{
			AuthService:  authService,
			IsValidToken: webApiService.IVk.Auth().TokenIsValid,
		},
		TelegramService: TelegramService{
			TGMenu:  menu.NewMenuService(webApiService, cfg.Buttons),
			SendMsg: webApiService.Send,
		},
	}
}
