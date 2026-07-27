package mocks

import "github.com/AppeiYA/consultation-platform/internal/identity/domain"

type MockSessionTokenGenerator struct {
	GenerateFn func() (domain.SessionToken, error)
}

func (m *MockSessionTokenGenerator) Generate() (domain.SessionToken, error) {
	if m.GenerateFn != nil {
		return m.GenerateFn()
	}
	return domain.SessionToken{}, nil
}