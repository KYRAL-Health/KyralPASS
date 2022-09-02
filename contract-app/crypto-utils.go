package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io"
)

func encrypt(key, text string) (string, error) {
	plaintext := []byte(text)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func decrypt(key, cryptotext, hash_ string) (string, bool, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(cryptotext)
	if err != nil {
		return "", false, err
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", false, err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", false, fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	hasher := sha512.New()
	hasher.Write(ciphertext)
	hashNew := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return string(ciphertext), hashNew == hash_, nil
}

func hash(text string) string {
	hasher := sha512.New()
	hasher.Write([]byte(text))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
