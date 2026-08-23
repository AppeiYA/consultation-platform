package outbound

import "context"

type RoleAssigner interface {
	AssignConsultantRole(ctx context.Context, userID string) error 
}