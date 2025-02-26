package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"

	"golang.org/x/crypto/bcrypt"
)

var (
	aesBlock     cipher.Block
	err          error
	gcmInstance  cipher.AEAD
	originalText []byte
	cipherText   []byte
)

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Encrypt(value []byte, keyPhrase string) ([]byte, error) {
	if aesBlock, err = aes.NewCipher([]byte(Md5Hashing(keyPhrase))); err != nil {
		return nil, err
	}

	if gcmInstance, err = cipher.NewGCM(aesBlock); err != nil {
		return nil, err
	}

	nonce := make([]byte, gcmInstance.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	cipherText = gcmInstance.Seal(nonce, nonce, value, nil)

	return cipherText, nil
}

func Decrypt(ciphered []byte, keyPhrase string) ([]byte, error) {
	if aesBlock, err = aes.NewCipher([]byte(Md5Hashing(keyPhrase))); err != nil {
		return nil, err
	}

	if gcmInstance, err = cipher.NewGCM(aesBlock); err != nil {
		return nil, err
	}

	nonceSize := gcmInstance.NonceSize()
	nonce, cipheredText := ciphered[:nonceSize], ciphered[nonceSize:]
	if originalText, err = gcmInstance.Open(nil, nonce, cipheredText, nil); err != nil {
		return nil, err
	}

	return originalText, nil
}

func Md5Hashing(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
