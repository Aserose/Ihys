package config

type Service struct {
	Telegram `yaml:"telegram"`
	Vk       `yaml:"vk"`
	Keypads  `yaml:"keypads"`
	LastFM
	Discogs
	Genius
	Auth
}

type Auth struct {
	Key string `env:"AES_KEY"`
}

type Telegram struct {
	WebhookLink string `yaml:"webhook_link"`
	Token       string `env:"TG_TOKEN"`
}

type Vk struct {
	AuthLink string `env:"AUTH_LINK"`
}

type LastFM struct {
	Key    string `env:"LASTFM_KEY"`
	Secret string `env:"LASTFM_SECRET"`
}

type Discogs struct {
	Key    string `env:"DISCOGS_KEY"`
	Secret string `env:"DISCOGS_SECRET"`
}

type Genius struct {
	Key    string `env:"GENIUS_KEY"`
	Secret string `env:"GENIUS_SECRET"`
}
