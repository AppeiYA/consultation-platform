package consultantIdentityAdapter

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/identity"
	identity_domain "github.com/AppeiYA/consultation-platform/internal/identity/domain"
)

type RoleAssigner struct {
	identityModule *identity.Module
}

func NewRoleAssigner(
    identityModule *identity.Module,
) *RoleAssigner {
    return &RoleAssigner{
        identityModule: identityModule,
    }
}

func (r *RoleAssigner) AssignConsultantRole(
    ctx context.Context,
    userID string,
) error {
    return r.identityModule.UpdateUserRole.Execute(
        ctx,
        userID,
        identity_domain.RoleConsultant,
    )
}