package repository

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/pkg/customLogger"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/jmoiron/sqlx"
	"path/filepath"
	"runtime"
)

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
