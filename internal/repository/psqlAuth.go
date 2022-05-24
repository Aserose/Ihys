package repository

import (
	"IhysBestowal/internal/dto"
	"IhysBestowal/pkg/customLogger"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type platforms struct {
	psqlVk
}

type psqlAuth struct {
	platforms
	nameTable string
	db        *sqlx.DB
	log       customLogger.Logger
}

func newPsqlAuth(log customLogger.Logger, db *sqlx.DB) psqlAuth {
	a := psqlAuth{
		platforms: platforms{
			psqlVk: newPsqlVk(log, db),
		},
		db:  db,
		log: log,
	}

	return a
}

func (a psqlAuth) Vk() IKey {
	return a.platforms.psqlVk
}

func (a psqlAuth) PutKey(user dto.TGUser, key string) {
	query := fmt.Sprintf(`INSERT INTO %s (tg_user_id, tg_chat_id, encrypted_key) VALUES ($1, $2, $3)`, a.nameTable)

	_, err := a.db.Query(query, user.UserId, user.ChatID, key)
	if err != nil {
		a.log.Error(a.log.CallInfoStr(), err.Error())
	}
}

func (a psqlAuth) GetKey(user dto.TGUser) (res string) {
	query := fmt.Sprintf(`SELECT encrypted_key FROM %s WHERE tg_user_id=$1 AND  tg_chat_id=$2`, a.nameTable)

	err := a.db.Get(&res, query, user.UserId, user.ChatID)
	if err != nil {
		a.log.Error(a.log.CallInfoStr(), err.Error())
	}

	return
}

func (a psqlAuth) HasKey(user dto.TGUser) bool {
	return a.GetKey(user) != ""
}

func (a psqlAuth) UpdateKey(user dto.TGUser, newKey string) {
	query := fmt.Sprintf(`UPDATE %s SET encrypted_key = $1 WHERE tg_user_id = $1 AND tg_chat_id = $2`, a.nameTable)

	_, err := a.db.Query(query, user.UserId, user.ChatID, newKey)
	if err != nil {
		a.log.Error(a.log.CallInfoStr(), err.Error())
	}
}

func (a psqlAuth) DeleteKey(user dto.TGUser) {
	query := fmt.Sprintf(`DELETE FROM %s WHERE tg_user_id = $1 AND tg_chat_id = $2`, a.nameTable)

	_, err := a.db.Query(query, user.UserId, user.ChatID)
	if err != nil {
		a.log.Error(a.log.CallInfoStr(), err.Error())
	}
}

type psqlVk struct{ IKey }

func newPsqlVk(log customLogger.Logger, db *sqlx.DB) psqlVk {
	return psqlVk{
		IKey: psqlAuth{
			db:        db,
			log:       log,
			nameTable: "vk",
		},
	}
}
