package domain

import (
	"time"

	custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"
)

const SessionIDPrefix = "sess"

type Session struct {
	id        string
	userID    string
	email string
	role string
	tokenHash string
	expiresAt time.Time
	createdAt time.Time
}

var (
	errEmptySessionID = custom_errors.BadException("session id is empty")
	errEmptyUserID    = custom_errors.BadException("user id is empty")
	errEmptyTokenHash = custom_errors.BadException("token hash is empty")
	errInvalidExpiry  = custom_errors.BadException("invalid expiry")
	errInvalidEmail = custom_errors.BadException("invalid email sent to session")
	errInvalidRole = custom_errors.BadException("invalid role")
)

func NewSession(
	id string,
	userID string,
	email string,
	role string,
	tokenHash string,
	now time.Time,
	expiry time.Duration,
) (*Session, error) {
	if id == "" {
		return nil, errEmptySessionID
	}
	if userID == "" {
		return nil, errEmptyUserID
	}
	address, err := NewEmail(email)
	if err != nil {
		return nil, errInvalidEmail
	}
	userRole, err := NewRole(role)
	if err != nil {
		return nil, errInvalidRole
	}
	if tokenHash == "" {
		return nil, errEmptyTokenHash
	}
	if expiry <= 0 {
		return nil, errInvalidExpiry
	}

	return &Session{
		id:        id,
		userID:    userID,
		email: address.String(),
		role: userRole.String(),
		tokenHash: tokenHash,
		createdAt: now,
		expiresAt: now.Add(expiry),
	}, nil
}

func (s *Session) ID() string {
	return s.id
}
func (s *Session) UserID() string {
	return s.userID
}
func (s *Session) Email() string {
	return s.email
}
func (s *Session) Role() string {
	return s.role
}
func (s *Session) TokenHash() string {
	return s.tokenHash
}
func (s *Session) CreatedAt() time.Time {
	return s.createdAt
}
func (s *Session) ExpiresAt() time.Time {
	return s.expiresAt
}

func (s Session) IsExpired(now time.Time) bool {
	return now.After(s.expiresAt)
}