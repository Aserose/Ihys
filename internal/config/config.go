package config

import (
	"IhysBestowal/pkg/customLogger"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"path/filepath"
	"runtime"
)

type Config struct {
	Service `yaml:"service"`
	Repository
	Server
	Handler `yaml:"handler"`
}

func NewCfg(log customLogger.Logger) (cfg Config) {
	err := godotenv.Load(getEnvPath()); if err != nil {
		log.Panic(log.CallInfoStr(), err.Error())
	}

	err = cleanenv.ReadEnv(&cfg); if err != nil {
		log.Panic(log.CallInfoStr(), err.Error())
	}

	err = cleanenv.ReadConfig(getYmlPath(), &cfg); if err != nil {
		log.Panic(log.CallInfoStr(), err.Error())
	}

	cfg.Service.Telegram.WebhookLink = cfg.Handler.Host + cfg.Handler.Api.Telegram

	return
}

func getYmlPath() string {
	return filepath.Join(getDirectory(getFilename(), 0), "config.yml")
}

func getEnvPath() string {
	return filepath.Join(getDirectory(getFilename(), 2), ".env")
}

func getDirectory(filename string, depth int) string {
	if depth != 0 { return getDirectory(filepath.Dir(filename), depth-1) }
	return filepath.Dir(filename)
}

func getFilename() string {
	_, filename, _, _ := runtime.Caller(0)
	return filename
}