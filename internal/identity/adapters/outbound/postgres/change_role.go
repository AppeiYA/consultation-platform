package postgres

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/identity/domain"
)

var updateUserRoleQuery = `
	UPDATE users SET role = $1 WHERE id = $2
`
func (a *UserRepository) ChangeRole(ctx context.Context, userID string, newRole domain.Role) error {
	executor := a.repository.Executor(ctx)
	_, err := executor.ExecContext(ctx, updateUserRoleQuery, newRole.String(), userID)
	if err != nil {
		return err
	}

	return nil
}