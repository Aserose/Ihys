package auth

import (
	"IhysBestowal/pkg/customLogger"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

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
