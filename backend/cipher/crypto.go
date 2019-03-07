package cipher

import (
	"crypto/md5"
	"encoding/hex"
)

type Cipher interface {
	Encrypt(data, key []byte) ([]byte, error)
	Decrypt(data, key []byte) ([]byte, error)
}

// hash converts a given slice of bytes to a fixed length one.
// This function is called by Encrypt and Decrypt for the key.
func hash(key []byte) []byte {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return []byte(hex.EncodeToString(hasher.Sum(nil)))
}
