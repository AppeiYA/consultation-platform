package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AppeiYA/consultation-platform/internal/identity/domain"
)

func(u *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	executor := u.repository.Executor(ctx)
	var model UserModel

	row := executor.QueryRowxContext(ctx, FIND_USER_BY_ID, id)
	if err := row.StructScan(&model); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return NewUserFromModel(model)
}