package mocks

import "context"

type MockRoleAssigner struct {
	AssignConsultantRoleFn func(ctx context.Context, userID string) error
}

func (m *MockRoleAssigner) AssignConsultantRole(ctx context.Context, userID string) error {
	if m.AssignConsultantRoleFn != nil {
		return m.AssignConsultantRoleFn(ctx, userID)
	}
	return nil
}
