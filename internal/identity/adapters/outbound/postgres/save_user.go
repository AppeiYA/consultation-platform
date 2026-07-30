package postgres

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/identity/domain"
)

func (u *UserRepository) Save(ctx context.Context, user *domain.User) error {
	executor := u.repository.Executor(ctx)
	model := NewUserModel(user)

	_, err := executor.NamedExecContext(ctx, CREATE_USER, model)
	return err
}