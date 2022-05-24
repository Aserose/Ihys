package repository

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/pkg/customLogger"
	"fmt"
	"github.com/jmoiron/sqlx"
)

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
