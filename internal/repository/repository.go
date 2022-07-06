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
	Close func()
}

func NewRepository(log customLogger.Logger, cfg config.Repository) Repository {
	psqls := newPsql(log, cfg.Psql)
	bdgr := newBadger(log)
	closeDB := func() {
		if err := bdgr.Close(); err != nil {
			log.Error(log.CallInfoStr(), err.Error())
		}
		if err := psqls.Close(); err != nil {
			log.Error(log.CallInfoStr(), err.Error())
		}
	}

	return Repository{
		Auth:         newPsqlAuth(log, psqls),
		TrackStorage: newBadgerTrackStorage(log, bdgr),
		Close:        closeDB,
	}
}

func newBadger(log customLogger.Logger) *badger.DB {
	_, filename, _, _ := runtime.Caller(0)
	badgerFilepath := filepath.Dir(filename) + `/badger/`

	bdgr, err := badger.Open(badger.DefaultOptions(badgerFilepath))
	if err != nil {
		log.Panic(log.CallInfoStr(), err.Error())
	}

	return bdgr
}

func newPsql(log customLogger.Logger, cfgPsql config.Psql) *sqlx.DB {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%s sslmode=%s",
		cfgPsql.Host, cfgPsql.DBName,
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
