package repository

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/dto"
	"IhysBestowal/pkg/customLogger"
	_ "github.com/lib/pq"
)

type Auth interface {
	Vk() Key
}

type Key interface {
	Create(usr dto.TGUser, key string)
	Get(usr dto.TGUser) string
	IsExist(usr dto.TGUser) bool
	Update(usr dto.TGUser, updKey string)
	Delete(usr dto.TGUser)
}

type Cache interface {
	Put(src datastruct.Song, similar datastruct.Set) string
	Get(src string, page int) []datastruct.Song
	IsExist(src string) bool
	PageCount(src string) int
	PageCapacity() int
}

type Repository struct {
	Auth
	Cache
	Close func()
}

func New(log customLogger.Logger, cfg config.Repository) Repository {
	psql := newPsql(log, cfg.Psql)
	bdgr := newBdgr(log)

	return Repository{
		Auth:  newPsqlAuth(log, psql),
		Cache: newBdgrCache(log, bdgr),
		Close: func() {
			if err := bdgr.Close(); err != nil {
				log.Error(log.CallInfoStr(), err.Error())
			}
			if err := psql.Close(); err != nil {
				log.Error(log.CallInfoStr(), err.Error())
			}
		},
	}
}
