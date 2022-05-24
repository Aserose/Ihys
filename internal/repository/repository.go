package repository

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/dto"
	"IhysBestowal/pkg/customLogger"
	_ "github.com/lib/pq"
)

type IKey interface {
	PutKey(user dto.TGUser, cryptKey string)
	GetKey(user dto.TGUser) string
	HasKey(user dto.TGUser) bool
	UpdateKey(user dto.TGUser, newCryptKey string)
	DeleteKey(user dto.TGUser)
}

type Auth interface {
	Vk() IKey
}

type Repository struct {
	Auth
}

func NewRepository(log customLogger.Logger, cfg config.Repository) Repository {
	psqls := newPsql(log, cfg.Psql)

	return Repository{
		Auth: newPsqlAuth(log, psqls),
	}
}
