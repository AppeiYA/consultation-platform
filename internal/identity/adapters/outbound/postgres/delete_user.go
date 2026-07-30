package postgres

import (
	"context"
)

func (u *UserRepository) Delete(ctx context.Context, id string) error {
	executor := u.repository.Executor(ctx)
	model := UserModel{
		ID:        id,
		IsDeleted: true,
		UpdatedAt: u.clock.Now(),
	}
	_, err := executor.NamedExecContext(ctx, DELETE_USER, model)
	return err
}