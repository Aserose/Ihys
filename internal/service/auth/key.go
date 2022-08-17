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

func (k Key) Create(usr dto.TGUser, key string) {
	enc, err := k.cypher.encrypt(key)
	if err != nil {
		k.log.Error(k.log.CallInfoStr(), err.Error())
	}

	k.key.Create(usr, enc)
}

func (k Key) Get(usr dto.TGUser) string {
	dec, err := k.cypher.decrypt(k.key.Get(usr))
	if err != nil {
		k.log.Error(k.log.CallInfoStr(), err.Error())
	}

	return dec
}

func (k Key) Update(usr dto.TGUser, newKey string) {
	enc, err := k.cypher.encrypt(newKey)
	if err != nil {
		k.log.Error(k.log.CallInfoStr(), err.Error())
	}

	k.key.Update(usr, enc)
}

func (k Key) IsExist(usr dto.TGUser) bool {
	return k.key.IsExist(usr)
}

func (k Key) Delete(usr dto.TGUser) {
	k.key.Delete(usr)
}
