package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strconv"

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

func NewIDCodec(key string) *IDCodec {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s", len(key), err)
		os.Exit(-1)
	}

	return &IDCodec{c}
}

func (c *IDCodec) Decid(b64 string) int {
	ciphertext, _ := base64.StdEncoding.DecodeString(b64)
	iv := ciphertext[:aes.BlockSize]
	plaintext := ciphertext[aes.BlockSize:]

	decrypter := cipher.NewCFBDecrypter(c.cipher, iv)
	decrypter.XORKeyStream(plaintext, plaintext)

	dec, _ := strconv.Atoi(string(plaintext))
	return dec
}

func (c *IDCodec) Encid(plainId int) string {
	plaintext := strconv.Itoa(plainId)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	encrypter := cipher.NewCFBEncrypter(c.cipher, iv)
	encrypter.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return base64.StdEncoding.EncodeToString(ciphertext)
}
