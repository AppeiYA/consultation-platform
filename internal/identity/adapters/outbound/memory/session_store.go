package memory

import (
	"context"
	"errors"

	"github.com/AppeiYA/consultation-platform/internal/identity/domain"
)

var ErrSessionNotFound = errors.New("session not found")

type SessionStore struct {
	sessions map[string]*domain.Session
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*domain.Session),
	}
}

func (s *SessionStore) Save(ctx context.Context, session *domain.Session) error {
	s.sessions[session.TokenHash()] = session
	return nil
}

func (s *SessionStore) FindByTokenHash(ctx context.Context, tokenHash string) (*domain.Session, error) {
	session, ok := s.sessions[tokenHash]
	if !ok {
		return nil, ErrSessionNotFound
	}
	return session, nil
}

func (s *SessionStore) Delete(ctx context.Context, tokenHash string) error {
	if _, ok := s.sessions[tokenHash]; !ok {
		return ErrSessionNotFound
	}
	delete(s.sessions, tokenHash)
	return nil
}