package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/identity/domain"
	"github.com/AppeiYA/consultation-platform/internal/identity/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/identity/usecase/dto"
)

type LoginUser struct {
	userRepository        outbound.UserRepository
	sessionStore          outbound.SessionStore
	passwordHasher        outbound.PasswordHasher
	sessionTokenHasher    outbound.SessionTokenHasher
	sessionTokenGenerator outbound.SessionTokenGenerator
	idGenerator           outbound.IdentifierGenerator
	clock                 outbound.Clock
}

func NewLoginUser(
	userRepository        outbound.UserRepository,
	sessionStore          outbound.SessionStore,
	passwordHasher        outbound.PasswordHasher,
	sessionTokenHasher    outbound.SessionTokenHasher,
	sessionTokenGenerator outbound.SessionTokenGenerator,
	idGenerator           outbound.IdentifierGenerator,
	clock                 outbound.Clock,
) *LoginUser {
	return &LoginUser{
		userRepository:        userRepository,
		sessionStore:          sessionStore,
		passwordHasher:        passwordHasher,
		sessionTokenHasher:    sessionTokenHasher,
		sessionTokenGenerator: sessionTokenGenerator,
		idGenerator:           idGenerator,
		clock:                 clock,
	}
}

func (l *LoginUser) Execute(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	email, err := domain.NewEmail(req.Email)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	password, err := domain.NewPassword(req.Password)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	user, err := l.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	ok, err := l.passwordHasher.Compare(password.String(), user.PasswordHash().String())
	if err != nil {
		return dto.LoginResponse{}, err
	}

	if !ok {
		return dto.LoginResponse{}, domain.ErrInvalidPassword
	}

	token, err := l.sessionTokenGenerator.Generate()
	if err != nil {
		return dto.LoginResponse{}, err
	}

	hash, err := l.sessionTokenHasher.Hash(token.String())
	if err != nil {
		return dto.LoginResponse{}, err
	}

	sessionID, err := l.idGenerator.Generate()
	if err != nil {
		return dto.LoginResponse{}, err
	}

	session, err := domain.NewSession(
		sessionID,
		user.ID(),
		hash,
		l.clock.Now(),
		7,
	)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	err = l.sessionStore.Save(ctx, session)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{
		SessionToken: token,
	}, nil
}