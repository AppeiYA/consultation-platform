package mocks 

import "errors"

type MockPasswordHasher struct {
	HashFn func(password string) (string, error)
	CompareFn func(hashedPassword, password string) (bool, error)
}

func (m *MockPasswordHasher) Hash(password string) (string, error) {
	if m.HashFn != nil {
		return m.HashFn(password)
	}
	return "", errors.New("HashFn not implemented")
}

func (m *MockPasswordHasher) Compare(hashedPassword, password string) (bool, error) {
	if m.CompareFn != nil {
		return m.CompareFn(hashedPassword, password)
	}
	return false, errors.New("CompareFn not implemented")
}