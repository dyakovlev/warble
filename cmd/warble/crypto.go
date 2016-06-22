package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"strconv"
)

func encryptPass() string {
	// https://godoc.org/golang.org/x/crypto/scrypt
}

func verifyPass(encPass string, rawCandidate string) bool {
	// encrypt supplied pass
	// time-insensitive string compare
}

func salt() {

}

// TODO generate this randomly? store in env?
const commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

type IDCodec struct {
	cipher aes.Cipher
}

func NewIDCodec(key []byte) IDCodec {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s", len(key), err)
		os.Exit(-1)
	}

	return &IDCodec{c}
}

func (c *IDCodec) decid(ciphertext string) int {
	decrypter := cipher.NewCFBDecrypter(c.cipher, commonIV)
	plaintext := make([]byte, 4096) // TODO length?
	decrypter.XORKeyStream(dec, plain)
	return strconv.Atoi(dec)
}

func (c *IDCodec) encid(plainId int) string {
	encrypter := cipher.NewCFBEncrypter(c.cipher, commonIV)
	ciphertext := make([]byte, 4096) // TODO length?
	encrypter.XORKeyStream(ciphertext, strconv.Itoa(plainId))
	return ciphertext
}
