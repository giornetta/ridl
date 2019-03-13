package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

type aesCrypto struct{}

// NewAES returns a concrete implementation of the Cipher interface.
// This implementation uses AES encryption.
func NewAES() Cipher {
	return &aesCrypto{}
}

// Encrypt encrypts a slice of bytes using a key.Encrypt
// This implementation uses AES
func (a *aesCrypto) Encrypt(data, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(hash(key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

// Decrypt tries to decrypt a given AES encrypted message using a key.
func (a *aesCrypto) Decrypt(data, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(hash(key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext is too short")
	}

	nonce, data := data[:nonceSize], data[nonceSize:]

	return gcm.Open(nil, nonce, data, nil)
}
