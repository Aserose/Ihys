package repository

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/dto"
	"IhysBestowal/pkg/customLogger"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"path/filepath"
	"runtime"
)

type Auth interface {
	Vk() IKey
}

type IKey interface {
	PutKey(user dto.TGUser, cryptKey string)
	GetKey(user dto.TGUser) string
	HasKey(user dto.TGUser) bool
	UpdateKey(user dto.TGUser, newCryptKey string)
	DeleteKey(user dto.TGUser)
}

type TrackStorage interface {
	Put(sourceAudio datastruct.AudioItem, similar datastruct.AudioItems) string
	GetItems(sourceAudio string, page int) []datastruct.AudioItem
	GetPageCount(sourceAudio string) int
	GetPageCapacity() int
	IsExist(sourceAudio string) bool
}

type Repository struct {
	Auth
	TrackStorage
}

func NewRepository(log customLogger.Logger, cfg config.Repository) Repository {
	psqls := newPsql(log, cfg.Psql)
	bdgr, err := newBadger()
	if err != nil {
		log.Panic(log.CallInfoStr(), err.Error())
	}

	return Repository{
		Auth:         newPsqlAuth(log, psqls),
		TrackStorage: newBadgerTrackStorage(log, bdgr),
	}
}

func newBadger() (*badger.DB, error) {
	_, filename, _, _ := runtime.Caller(0)
	return badger.Open(badger.DefaultOptions(filepath.Dir(filename) + `/badger/storage`))
}

func newPsql(log customLogger.Logger, cfgPsql config.Psql) *sqlx.DB {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s port=%s sslmode=%s",
		cfgPsql.User, cfgPsql.Password,
		cfgPsql.Port, cfgPsql.SSLMode))
	if err != nil {
		log.Panic(log.CallInfoStr(), err.Error())
	}

	if err = db.Ping(); err != nil {
		log.Panic(log.CallInfoStr(), err.Error())
	}

	return db
}
