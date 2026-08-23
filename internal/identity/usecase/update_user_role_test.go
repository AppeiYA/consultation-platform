package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/AppeiYA/consultation-platform/internal/identity/domain"
	"github.com/AppeiYA/consultation-platform/internal/identity/mocks"
)

func TestUpdateUserRole_Execute(t *testing.T) {
	userID := "user_123"
	newRole := domain.RoleConsultant

	t.Run("should update user role successfully", func(t *testing.T) {
		var capturedUserID string
		var capturedRole domain.Role

		repo := &mocks.MockUserRepository{
			ChangeRoleFn: func(ctx context.Context, uID string, role domain.Role) error {
				capturedUserID = uID
				capturedRole = role
				return nil
			},
		}

		uc := NewUpdateUserRoleUsecase(repo)
		err := uc.Execute(context.Background(), userID, newRole)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if capturedUserID != userID {
			t.Errorf("expected userID %s, got %s", userID, capturedUserID)
		}
		if capturedRole != newRole {
			t.Errorf("expected role %s, got %s", newRole, capturedRole)
		}
	})

	t.Run("should fail when user repository ChangeRole returns error", func(t *testing.T) {
		expectedErr := errors.New("db error")
		repo := &mocks.MockUserRepository{
			ChangeRoleFn: func(ctx context.Context, uID string, role domain.Role) error {
				return expectedErr
			},
		}

		uc := NewUpdateUserRoleUsecase(repo)
		err := uc.Execute(context.Background(), userID, newRole)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})
}
