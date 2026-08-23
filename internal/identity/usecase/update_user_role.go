package usecase

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/identity/domain"
	"github.com/AppeiYA/consultation-platform/internal/identity/ports/outbound"
)

type UpdateUserRoleUsecase struct {
	userRepository outbound.UserRepository
}

func NewUpdateUserRoleUsecase(
	userRepository outbound.UserRepository,
) *UpdateUserRoleUsecase {
	return &UpdateUserRoleUsecase{
		userRepository: userRepository,
	}
}

func (uc *UpdateUserRoleUsecase) Execute(ctx context.Context, userID string, newRole domain.Role) error {
	return uc.userRepository.ChangeRole(ctx, userID, newRole)
}