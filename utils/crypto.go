package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"strconv"

	// despite being called simple-scrypt this provides an "scrypt" package
	"github.com/elithrar/simple-scrypt"
)

func EncryptPass(rawPassword string) string {
	hash, _ := scrypt.GenerateFromPassword([]byte(rawPassword), scrypt.DefaultParams)
	return string(hash)
}

func VerifyPass(encPass string, rawCandidate string) bool {
	err := scrypt.CompareHashAndPassword([]byte(encPass), []byte(rawCandidate))
	return err == nil
}

type IDCodec struct {
	cipher cipher.Block
}

func NewIDCodec(key string) (*IDCodec, error) {
	c, err := aes.NewCipher([]byte(key))
	return &IDCodec{c}, err
}

func (c *IDCodec) Decid(b64 string) int64 {
	ciphertext, _ := base64.StdEncoding.DecodeString(b64)
	iv := ciphertext[:aes.BlockSize]
	plaintext := ciphertext[aes.BlockSize:]

	decrypter := cipher.NewCFBDecrypter(c.cipher, iv)
	decrypter.XORKeyStream(plaintext, plaintext)

	dec, _ := strconv.ParseInt(string(plaintext), 10, 64)
	return dec
}

func (c *IDCodec) Encid(plainId int64) string {
	plaintext := strconv.FormatInt(plainId, 10)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	encrypter := cipher.NewCFBEncrypter(c.cipher, iv)
	encrypter.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return base64.StdEncoding.EncodeToString(ciphertext)
}
