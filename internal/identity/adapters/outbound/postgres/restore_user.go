package postgres

import (
	"context"
)

func (u *UserRepository) Restore(ctx context.Context, id string) error {
	executor := u.repository.Executor(ctx)
	model := UserModel{
		ID:        id,
		IsDeleted: false,
		UpdatedAt: u.clock.Now(),
	}
	_, err := executor.NamedExecContext(ctx, RESTORE_USER, model)
	return err
}
