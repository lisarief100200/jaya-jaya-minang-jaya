package helpers

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hashed(value string) string {
	valueBytes := []byte(value)
	sha256Hasher := sha256.New()
	sha256Hasher.Write(valueBytes)
	hashedValueBytes := sha256Hasher.Sum(nil)
	hashedValueHex := hex.EncodeToString(hashedValueBytes)

	return hashedValueHex
}
