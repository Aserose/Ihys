package auth

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/repository"
	"IhysBestowal/pkg/customLogger"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

type IKey interface {
	PutKey(user dto.TGUser, key string)
	GetKey(user dto.TGUser) string
	HasKey(user dto.TGUser) bool
	UpdateKey(user dto.TGUser, newKey string)
	DeleteKey(user dto.TGUser)
}

type AuthService interface {
	Vk() IKey
}

type platforms struct {
	authVk
}

type authService struct {
	platforms
	cypher
	repo repository.Repository
	key  func() repository.IKey
	log  customLogger.Logger
}

func NewAuthService(log customLogger.Logger, cfg config.Auth, repo repository.Repository) AuthService {
	cyph := newCypher(log, []byte(cfg.Key))

	return authService{
		platforms: platforms{
			authVk: newAuthVk(log, repo, cyph),
		},
		cypher: cyph,
		repo:   repo,
		log:    log,
	}
}

func (as authService) Vk() IKey {
	return as.platforms.authVk
}

func (as authService) PutKey(user dto.TGUser, key string) {
	encryptedKey, err := as.cypher.encrypt(key)
	if err != nil {
		as.log.Error(as.log.CallInfoStr(), err.Error())
	}

	as.key().PutKey(user, encryptedKey)
}

func (as authService) GetKey(user dto.TGUser) string {
	decryptedKey, err := as.cypher.decrypt(as.key().GetKey(user))
	if err != nil {
		as.log.Error(as.log.CallInfoStr(), err.Error())
	}

	return decryptedKey
}

func (as authService) UpdateKey(user dto.TGUser, newKey string) {
	encryptedKey, err := as.cypher.encrypt(newKey)
	if err != nil {
		as.log.Error(as.log.CallInfoStr(), err.Error())
	}

	as.key().UpdateKey(user, encryptedKey)
}

func (as authService) HasKey(user dto.TGUser) bool {
	return as.key().HasKey(user)
}

func (as authService) DeleteKey(user dto.TGUser) {
	as.key().DeleteKey(user)
}

type cypher struct {
	key []byte
	log customLogger.Logger
}

func newCypher(log customLogger.Logger, key []byte) cypher {
	return cypher{
		key: key,
		log: log,
	}
}

func (c cypher) encrypt(key string) (string, error) {
	byteMsg := []byte(key)
	block, err := aes.NewCipher(c.key)
	if err != nil {
		c.log.Error(c.log.CallInfoStr(), err.Error())
	}

	cipherText := make([]byte, aes.BlockSize+len(byteMsg))
	iv := cipherText[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		c.log.Error(c.log.CallInfoStr(), err.Error())
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], byteMsg)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (c cypher) decrypt(key string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		c.log.Error(c.log.CallInfoStr(), err.Error())
	}

	block, err := aes.NewCipher(c.key)
	if err != nil {
		c.log.Error(c.log.CallInfoStr(), err.Error())
	}

	if len(cipherText) < aes.BlockSize {
		return "", fmt.Errorf("invalid ciphertext block size")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}

type authVk struct {
	IKey
}

func newAuthVk(log customLogger.Logger, repo repository.Repository, cyph cypher) authVk {
	return authVk{
		authService{
			cypher: cyph,
			repo:   repo,
			log:    log,
			key:    repo.Vk,
		},
	}
}
