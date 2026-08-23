package inbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/identity/domain"
)

type UpdateUserRoleInt interface {
	Execute(ctx context.Context, userID string, newRole domain.Role) error
}