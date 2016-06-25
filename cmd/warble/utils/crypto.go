package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"
	"strconv"

	"github.com/elithrar/simple-scrypt"
)

func EncryptPass(rawPassword string) string {
	hash, err := scrypt.GenerateFromPassword([]byte(rawPassword), scrypt.DefaultParams)
	return string(hash)
}

func VerifyPass(encPass string, rawCandidate string) bool {
	err := scrypt.CompareHashAndPassword([]byte(encPass), []byte(rawCandidate))
	return err == nil
}

// TODO generate this randomly per encryption, cc. https://gist.github.com/manishtpatel/8222606
var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

type IDCodec struct {
	cipher cipher.Block
}

func NewIDCodec(key string) *IDCodec {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s", len(key), err)
		os.Exit(-1)
	}

	return &IDCodec{c}
}

func (c *IDCodec) Decid(ciphertext string) int {
	decrypter := cipher.NewCFBDecrypter(c.cipher, commonIV)
	plaintext := make([]byte, 4096) // TODO length?
	decrypter.XORKeyStream([]byte(ciphertext), plaintext)
	dec, err := strconv.Atoi(string(plaintext))
	return dec
}

func (c *IDCodec) Encid(plainId int) string {
	encrypter := cipher.NewCFBEncrypter(c.cipher, commonIV)
	ciphertext := make([]byte, 4096) // TODO length?
	encrypter.XORKeyStream(ciphertext, []byte(strconv.Itoa(plainId)))
	return string(ciphertext)
}
