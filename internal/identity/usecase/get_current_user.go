package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/identity/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/identity/usecase/dto"
)

type GetCurrentUser struct {
	userRepository outbound.UserRepository
}

func NewGetCurrentUser(userRepository outbound.UserRepository) *GetCurrentUser {
	return &GetCurrentUser{userRepository: userRepository}
}

func (g *GetCurrentUser) Execute(ctx context.Context, req dto.GetCurrentUserRequest) (dto.GetCurrentUserResponse, error) {
	user, err := g.userRepository.FindByID(ctx, req.UserID)
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