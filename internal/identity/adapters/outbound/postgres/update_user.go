package postgres

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/identity/domain"
)

func (u *UserRepository) Update(ctx context.Context, user *domain.User) error {
	executor := u.repository.Executor(ctx)
	model := NewUserModel(user)

	_, err := executor.NamedExecContext(ctx, UPDATE_USER, model)
	return err	
}