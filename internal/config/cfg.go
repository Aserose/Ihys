package config

import (
	"IhysBestowal/pkg/customLogger"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"path/filepath"
	"runtime"
)

type Cfg struct {
	Service `yaml:"service"`
	Repository
	Server
	Handler `yaml:"handler"`
}

func New(log customLogger.Logger) (cfg Cfg) {
	err := godotenv.Load(envPath())
	if err != nil {
		log.Panic(log.CallInfoStr(), err.Error())
	}

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Panic(log.CallInfoStr(), err.Error())
	}

	err = cleanenv.ReadConfig(ymlPath(), &cfg)
	if err != nil {
		log.Panic(log.CallInfoStr(), err.Error())
	}

	cfg.Service.Telegram.WebhookLink = cfg.Handler.Host + cfg.Handler.Api.Telegram

	return
}

func ymlPath() string {
	return filepath.Join(directory(filename(), 0), "config.yml")
}

func envPath() string {
	return filepath.Join(directory(filename(), 2), ".env")
}

func directory(filename string, depth int) string {
	if depth != 0 {
		return directory(filepath.Dir(filename), depth-1)
	}
	return filepath.Dir(filename)
}

func filename() string {
	_, filename, _, _ := runtime.Caller(0)
	return filename
}
