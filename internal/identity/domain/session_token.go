package domain

import (
	custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"
)

type SessionToken struct {
	value string
}

var (
	errEmptySessionToken         = custom_errors.BadException("session token is empty")
	errSessionTokenTooShort      = custom_errors.BadException("session token is too short")
)

const (
	minSessionTokenLength = 32
)

func NewSessionToken(value string) (SessionToken, error) {
	if value == "" {
		return SessionToken{}, errEmptySessionToken
	}

	if len(value) < minSessionTokenLength {
		return SessionToken{}, errSessionTokenTooShort
	}

	return SessionToken{value: value}, nil
}

func (s SessionToken) String() string {
	return s.value
}

func (s SessionToken) IsZero() bool {
	return s.value == ""
}