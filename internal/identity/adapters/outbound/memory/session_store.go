package memory

import (
	"context"
	"sync"

	"github.com/AppeiYA/consultation-platform/internal/identity/domain"
	custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"
)

var ErrSessionNotFound = custom_errors.InternalServerError("session not found")

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*domain.Session
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*domain.Session),
	}
}

func (s *SessionStore) Save(ctx context.Context, session *domain.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[session.TokenHash()] = session
	return nil
}

func (s *SessionStore) FindByTokenHash(ctx context.Context, tokenHash string) (*domain.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, ok := s.sessions[tokenHash]
	if !ok {
		return nil, ErrSessionNotFound
	}
	return session, nil
}

func (s *SessionStore) Delete(ctx context.Context, tokenHash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.sessions[tokenHash]; !ok {
		return ErrSessionNotFound
	}
	delete(s.sessions, tokenHash)
	return nil
}