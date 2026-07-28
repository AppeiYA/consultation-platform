package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

const defaultCost = bcrypt.DefaultCost // 10

type Hasher struct {
	cost int
}

func NewHasher() *Hasher {
	return &Hasher{cost: defaultCost}
}

func (h *Hasher) Hash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (h *Hasher) Compare(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err
	}

	return true, nil
}