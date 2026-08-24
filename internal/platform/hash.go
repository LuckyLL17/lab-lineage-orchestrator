package platform

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func Hash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func ChainHash(previous, value string) string {
	return Hash(strings.Join([]string{previous, value}, "|"))
}
