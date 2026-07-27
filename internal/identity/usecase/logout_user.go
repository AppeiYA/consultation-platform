package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/identity/domain"
	"github.com/AppeiYA/consultation-platform/internal/identity/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/identity/usecase/dto"
)

type LogoutUser struct {
	sessionStore outbound.SessionStore
	sessionTokenHasher outbound.SessionTokenHasher
}

func NewLogoutUser(
	sessionStore outbound.SessionStore,
	sessionTokenHasher outbound.SessionTokenHasher,
) *LogoutUser {
	return &LogoutUser{
		sessionStore:          sessionStore,
		sessionTokenHasher:    sessionTokenHasher,
	}
}

func (l *LogoutUser) Execute(ctx context.Context, req dto.LogoutRequest) (dto.LogoutResponse, error) {
	token, err := domain.NewSessionToken(req.SessionToken)
	if err != nil {
		return dto.LogoutResponse{}, err
	}

	hash, err := l.sessionTokenHasher.Hash(token.String())
	if err != nil {
		return dto.LogoutResponse{}, err
	}

	err = l.sessionStore.Delete(ctx, hash)
	if err != nil {
		return dto.LogoutResponse{}, err
	}

	return dto.LogoutResponse{}, nil 
}