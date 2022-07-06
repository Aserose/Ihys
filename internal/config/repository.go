package config

type Repository struct {
	Psql
}

type Psql struct {
	Host     string `env:"PSQL_HOST"`
	DBName   string `env:"PSQL_NAME"`
	User     string `env:"PSQL_USER"`
	Password string `env:"PSQL_PASSWORD"`
	Port     string `env:"PSQL_PORT"`
	SSLMode  string `env:"PSQL_SSLMODE"`
}
