package auth

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/repository"
	"IhysBestowal/pkg/customLogger"
)

type Key struct {
	cypher
	key repository.Key
}

func newKey(log customLogger.Logger, cfg config.Auth, repo repository.Key) Key {
	return Key{
		cypher: newCypher(log, []byte(cfg.Key)),
		key:    repo,
	}
}

func (as Key) Create(user dto.TGUser, key string) {
	encryptedKey, err := as.cypher.encrypt(key)
	if err != nil {
		as.log.Error(as.log.CallInfoStr(), err.Error())
	}

	as.key.Create(user, encryptedKey)
}

func (as Key) Get(user dto.TGUser) string {
	decryptedKey, err := as.cypher.decrypt(as.key.Get(user))
	if err != nil {
		as.log.Error(as.log.CallInfoStr(), err.Error())
	}

	return decryptedKey
}

func (as Key) Update(user dto.TGUser, newKey string) {
	encryptedKey, err := as.cypher.encrypt(newKey)
	if err != nil {
		as.log.Error(as.log.CallInfoStr(), err.Error())
	}

	as.key.Update(user, encryptedKey)
}

func (as Key) IsExist(user dto.TGUser) bool {
	return as.key.IsExist(user)
}

func (as Key) Delete(user dto.TGUser) {
	as.key.Delete(user)
}
