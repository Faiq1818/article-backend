package pkg

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func RandomHash() (string, error) {
	randomBytes := make([]byte, 32)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(randomBytes)

	return string(hex.EncodeToString(hash[:])), nil
}
