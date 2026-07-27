package mocks

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/identity/domain"
)

type MockSessionStore struct {
	SaveFn            func(ctx context.Context, session *domain.Session) error
	FindByTokenHashFn func(ctx context.Context, tokenHash string) (*domain.Session, error)
	DeleteFn          func(ctx context.Context, sessionID string) error
}

func (m *MockSessionStore) Save(ctx context.Context, session *domain.Session) error {
	if m.SaveFn != nil {
		return m.SaveFn(ctx, session)
	}
	return nil
}

func (m *MockSessionStore) FindByTokenHash(ctx context.Context, tokenHash string) (*domain.Session, error) {
	if m.FindByTokenHashFn != nil {
		return m.FindByTokenHashFn(ctx, tokenHash)
	}
	return nil, domain.ErrSessionNotFound
}

func (m *MockSessionStore) Delete(ctx context.Context, sessionID string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, sessionID)
	}
	return nil
}