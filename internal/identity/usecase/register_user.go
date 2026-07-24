package usecase

import (
	"context"
	"errors"

	"github.com/AppeiYA/consultation-platform/internal/identity/domain"
	"github.com/AppeiYA/consultation-platform/internal/identity/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/identity/usecase/dto"
)

type RegisterUser struct {
	userRepo outbound.UserRepository
	passwordHasher outbound.PasswordHasher
	idGenerator outbound.IdentifierGenerator
	clock outbound.Clock
}

type RegisterUserParams struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	Role      string
}

func NewRegisterUser(
	userRepo outbound.UserRepository, 
	passwordHasher outbound.PasswordHasher, 
	idGenerator outbound.IdentifierGenerator, 
	clock outbound.Clock,
	) *RegisterUser {
	return &RegisterUser{
		userRepo: userRepo, 
		passwordHasher: passwordHasher, 
		idGenerator: idGenerator,
		clock: clock,
	}
}

func (r *RegisterUser) Execute(ctx context.Context, params RegisterUserParams) (*dto.RegisterUserResponse, error) {
	email, err := domain.NewEmail(params.Email)
	if err != nil {
		return nil, err
	}

	password, err := domain.NewPassword(params.Password)
	if err != nil {
		return nil, err
	}

	role, err := domain.NewRole(params.Role)
	if err != nil {
		return nil, err
	}

	exists, err := r.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			return nil, err
		}
	}

	if exists != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	hashedPassword, err := r.passwordHasher.Hash(password.String())
	if err != nil {
		return nil, err
	}

	userID, err := r.idGenerator.Generate()
	if err != nil {
		return nil, err
	}

	passwordHash := domain.NewPasswordHash(hashedPassword)

	user := domain.NewUser(
		userID,
		params.FirstName,
		params.LastName,
		email,
		passwordHash,
		role,
		r.clock.Now(),
	)

	err = r.userRepo.Save(ctx, user)
	if err != nil {
		return nil, err
	}

	return &dto.RegisterUserResponse{
		UserID: user.ID(),
	}, nil
}