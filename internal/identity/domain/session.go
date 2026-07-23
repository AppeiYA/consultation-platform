package domain

import (
	"time"

	custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"
)

type Session struct {
	id        string
	userID    string
	tokenHash string
	expiresAt time.Time
	createdAt time.Time
}

var (
	errEmptySessionID = custom_errors.BadException("session id is empty")
	errEmptyUserID    = custom_errors.BadException("user id is empty")
	errEmptyTokenHash = custom_errors.BadException("token hash is empty")
	errInvalidExpiry  = custom_errors.BadException("invalid expiry")
)

func NewSession(
	id string,
	userID string,
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
	if tokenHash == "" {
		return nil, errEmptyTokenHash
	}
	if expiry <= 0 {
		return nil, errInvalidExpiry
	}

	return &Session{
		id:        id,
		userID:    userID,
		tokenHash: tokenHash,
		createdAt: now,
		expiresAt: now.Add(expiry),
	}, nil
}

func (s *Session) ID() string
func (s *Session) UserID() string
func (s *Session) TokenHash() string
func (s *Session) CreatedAt() time.Time
func (s *Session) ExpiresAt() time.Time

func (s Session) IsExpired(now time.Time) bool {
	return now.After(s.expiresAt)
}