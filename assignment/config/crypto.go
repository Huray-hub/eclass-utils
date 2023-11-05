package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

func encrypt(plaintext, secretKey string) (ciphertext string, err error) {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure random sequence
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return
	}

	ciphertextBytes := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	ciphertext = string(ciphertextBytes)
	return
}

func decrypt(ciphertext, secretKey string) (plaintext string, err error) {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		err = errors.New("decryption: nonce's character length is bigger than ciphetext's")
		return
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintextBytes, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", err
	}

	plaintext = string(plaintextBytes)
	return
}

func generateSecretKey() (string, error) {
	key := make([]byte, 32)

	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return string(key), nil
}
