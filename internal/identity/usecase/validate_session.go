package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/identity/domain"
	"github.com/AppeiYA/consultation-platform/internal/identity/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/identity/usecase/dto"
)

type ValidateSession struct {
	sessionStore       outbound.SessionStore
	sessionTokenHasher outbound.SessionTokenHasher
	clock              outbound.Clock
}

func NewValidateSession(sessionStore outbound.SessionStore, sessionTokenHasher outbound.SessionTokenHasher, clock outbound.Clock) *ValidateSession {
	return &ValidateSession{
		sessionStore:       sessionStore,
		sessionTokenHasher: sessionTokenHasher,
		clock:              clock,
	}
}

func (uc *ValidateSession) Execute(ctx context.Context, req dto.ValidateSessionRequest) (dto.ValidateSessionResponse, error) {
	token, err := domain.NewSessionToken(req.SessionToken)
	if err != nil {
		return dto.ValidateSessionResponse{}, err
	}

	hashed, err := uc.sessionTokenHasher.Hash(token.String())
	if err != nil {
		return dto.ValidateSessionResponse{}, err
	}

	session, err := uc.sessionStore.FindByTokenHash(ctx, hashed)
	if err != nil {
		return dto.ValidateSessionResponse{}, err
	}

	if session.IsExpired(uc.clock.Now()) {
		return dto.ValidateSessionResponse{}, domain.ErrSessionExpired
	}

	return dto.ValidateSessionResponse{
		SessionClaims: dto.SessionClaims{
			UserID: session.UserID(),
			Email:  session.Email(),
			Role:   session.Role(),
		},
	}, nil
}