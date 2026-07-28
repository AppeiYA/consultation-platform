package random

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/AppeiYA/consultation-platform/internal/identity/domain"
)

const tokenBytes = 32 // 256 bits of entropy

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate() (domain.SessionToken, error) {
	b := make([]byte, tokenBytes)
	if _, err := rand.Read(b); err != nil {
		return domain.SessionToken{}, err
	}
	return domain.NewSessionToken(base64.RawURLEncoding.EncodeToString(b))
}