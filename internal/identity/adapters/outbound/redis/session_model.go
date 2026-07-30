package redis

import (
	"time"

	"github.com/AppeiYA/consultation-platform/internal/identity/domain"
)

type SessionModel struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	TokenHash string    `json:"token_hash"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}


func NewSessionModel(session *domain.Session) *SessionModel {
	return &SessionModel{
		ID:        session.ID(),
		UserID:    session.UserID(),
		Email:     session.Email(),
		Role:      session.Role(),
		TokenHash: session.TokenHash(),
		CreatedAt: session.CreatedAt(),
		ExpiresAt: session.ExpiresAt(),
	}
}

func (m *SessionModel) ToDomain() (*domain.Session, error) {
	return domain.ReconstituteSession(
		m.ID,
		m.UserID,
		m.Email,
		m.Role,
		m.TokenHash,
		m.CreatedAt,
		m.ExpiresAt,
	)
}