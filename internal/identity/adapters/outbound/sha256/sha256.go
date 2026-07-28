package sha256

import (
	"crypto/sha256"
	"encoding/hex"
)

type Hasher struct{}

func NewHasher() *Hasher {
	return &Hasher{}
}

func (h *Hasher) Hash(token string) (string, error) {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:]), nil
}