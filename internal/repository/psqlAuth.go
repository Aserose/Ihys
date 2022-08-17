package repository

import (
	"IhysBestowal/internal/dto"
	"IhysBestowal/pkg/customLogger"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type auth struct {
	nameTable string
	db        *sqlx.DB
	log       customLogger.Logger
}

type psqlAuth struct {
	vk auth
}

func newPsqlAuth(log customLogger.Logger, db *sqlx.DB) psqlAuth {
	return psqlAuth{
		vk: auth{
			db:        db,
			log:       log,
			nameTable: "vk",
		},
	}
}

func (a psqlAuth) Vk() Key { return a.vk }

func (a auth) Create(usr dto.TGUser, key string) {
	query := fmt.Sprintf(`INSERT INTO %s (tg_user_id, tg_chat_id, encrypted_key) VALUES ($1, $2, $3)`, a.nameTable)

	_, err := a.db.Query(query, usr.UserId, usr.ChatId, key)
	if err != nil {
		a.log.Error(a.log.CallInfoStr(), err.Error())
	}
}

func (a auth) Get(usr dto.TGUser) (res string) {
	query := fmt.Sprintf(`SELECT encrypted_key FROM %s WHERE tg_user_id=$1 AND  tg_chat_id=$2`, a.nameTable)

	err := a.db.Get(&res, query, usr.UserId, usr.ChatId)
	if err != nil {
		a.log.Error(a.log.CallInfoStr(), err.Error())
	}

	return
}

func (a auth) IsExist(usr dto.TGUser) bool {
	return a.Get(usr) != ""
}

func (a auth) Update(usr dto.TGUser, newKey string) {
	query := fmt.Sprintf(`UPDATE %s SET encrypted_key = $1 WHERE tg_user_id = $1 AND tg_chat_id = $2`, a.nameTable)

	_, err := a.db.Query(query, usr.UserId, usr.ChatId, newKey)
	if err != nil {
		a.log.Error(a.log.CallInfoStr(), err.Error())
	}
}

func (a auth) Delete(usr dto.TGUser) {
	query := fmt.Sprintf(`DELETE FROM %s WHERE tg_user_id = $1 AND tg_chat_id = $2`, a.nameTable)

	_, err := a.db.Query(query, usr.UserId, usr.ChatId)
	if err != nil {
		a.log.Error(a.log.CallInfoStr(), err.Error())
	}
}
