package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/identity/domain"
	"github.com/AppeiYA/consultation-platform/internal/identity/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/identity/usecase/dto"
)

type GetCurrentUser struct {
	sessionStore       outbound.SessionStore
	sessionTokenHasher outbound.SessionTokenHasher
	userRepository     outbound.UserRepository
	clock              outbound.Clock
}

func NewGetCurrentUser(
	sessionStore outbound.SessionStore,
	sessionTokenHasher outbound.SessionTokenHasher,
	userRepository     outbound.UserRepository,
	clock              outbound.Clock,
) *GetCurrentUser {
	return &GetCurrentUser{
		sessionStore:       sessionStore,
		sessionTokenHasher: sessionTokenHasher,
		userRepository:     userRepository,
		clock:              clock,
	}
}

func (g *GetCurrentUser) Execute(ctx context.Context, req dto.GetCurrentUserRequest) (dto.GetCurrentUserResponse, error) {
	token, err := domain.NewSessionToken(req.SessionToken)
	if err != nil {
		return dto.GetCurrentUserResponse{}, err
	}

	hashed, err := g.sessionTokenHasher.Hash(token.String())
	if err != nil {
		return dto.GetCurrentUserResponse{}, err
	}

	session, err := g.sessionStore.FindByTokenHash(ctx, hashed)
	if err != nil {
		return dto.GetCurrentUserResponse{}, err
	}

	if session.IsExpired(g.clock.Now()) {
		return dto.GetCurrentUserResponse{}, domain.ErrSessionExpired
	}

	user, err := g.userRepository.FindByID(ctx, session.UserID())
	if err != nil {
		return dto.GetCurrentUserResponse{}, err
	}

	return dto.GetCurrentUserResponse{
		ID:        user.ID(),
		FirstName: user.FirstName(),
		LastName:  user.LastName(),
		Email:     user.Email().String(),
		Role:      user.Role().String(),
	}, nil
}