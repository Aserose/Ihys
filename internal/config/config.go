package config

import (
	"IhysBestowal/pkg/customLogger"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Service `yaml:"service"`
	Repository
	Server
	Handler `yaml:"handler"`
}

func NewCfg(log customLogger.Logger) (cfg Config) {
	err := godotenv.Load(); if err != nil {
		log.Panic(log.CallInfoStr(), err.Error())
	}

	err = cleanenv.ReadEnv(&cfg); if err != nil {
		log.Panic(log.CallInfoStr(), err.Error())
	}

	err = cleanenv.ReadConfig("./internal/config/config.yml", &cfg); if err != nil {
		log.Panic(log.CallInfoStr(), err.Error())
	}

	cfg.Service.Telegram.WebhookLink = cfg.Handler.Host + cfg.Handler.Api.Telegram

	return
}
