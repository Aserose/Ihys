package config

type Repository struct {
	Psql
}

type Psql struct {
	User     string `env:"PSQL_USER"`
	Password string `env:"PSQL_PASSWORD"`
	Port     string `env:"PSQL_PORT"`
	SSLMode  string `env:"PSQL_SSLMODE"`
}
