package postgres

import (
	"context"
)

func (u *UserRepository) Verify(ctx context.Context, id string) error {
	executor := u.repository.Executor(ctx)
	model := UserModel{
		ID:        id,
		UpdatedAt: u.clock.Now(),
	}
	_, err := executor.NamedExecContext(ctx, VERIFY_USER, model)
	return err
}
