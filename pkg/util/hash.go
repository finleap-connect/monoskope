package util

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashString(input string) string {
	return HashBytes([]byte(input))
}

func HashBytes(input []byte) string {
	sha256Bytes := sha256.Sum256(input)
	return hex.EncodeToString(sha256Bytes[:])
}
