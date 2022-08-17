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

func (k Key) Create(user dto.TGUser, key string) {
	enc, err := k.cypher.encrypt(key)
	if err != nil {
		k.log.Error(k.log.CallInfoStr(), err.Error())
	}

	k.key.Create(user, enc)
}

func (k Key) Get(user dto.TGUser) string {
	dec, err := k.cypher.decrypt(k.key.Get(user))
	if err != nil {
		k.log.Error(k.log.CallInfoStr(), err.Error())
	}

	return dec
}

func (k Key) Update(user dto.TGUser, newKey string) {
	enc, err := k.cypher.encrypt(newKey)
	if err != nil {
		k.log.Error(k.log.CallInfoStr(), err.Error())
	}

	k.key.Update(user, enc)
}

func (k Key) IsExist(user dto.TGUser) bool {
	return k.key.IsExist(user)
}

func (k Key) Delete(user dto.TGUser) {
	k.key.Delete(user)
}
