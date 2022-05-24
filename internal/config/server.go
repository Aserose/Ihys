package config

type Server struct {
	Port string `env:"SERVER_PORT"`
}

type Handler struct {
	Host string `env:"HOST"`
	Api  `yaml:"api"`
}

type Api struct {
	Telegram string `yaml:"telegram"`
}
