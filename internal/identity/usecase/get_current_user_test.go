package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/identity/domain"
	"github.com/AppeiYA/consultation-platform/internal/identity/mocks"
	"github.com/AppeiYA/consultation-platform/internal/identity/usecase/dto"
)

type testGetCurrentUser struct {
	sut            *GetCurrentUser
	userRepository *mocks.MockUserRepository
}

func setupGetCurrentUser(t *testing.T) *testGetCurrentUser {
	t.Helper()

	userRepo := &mocks.MockUserRepository{
		FindByIDFn: func(ctx context.Context, id string) (*domain.User, error) {
			return nil, domain.ErrUserNotFound
		},
		FindByEmailFn: func(ctx context.Context, email domain.Email) (*domain.User, error) {
			return nil, domain.ErrUserNotFound
		},
		SaveFn: func(ctx context.Context, user *domain.User) error {
			return nil
		},
	}

	sut := NewGetCurrentUser(userRepo)

	return &testGetCurrentUser{
		sut:            sut,
		userRepository: userRepo,
	}
}

func TestGetCurrentUser_Execute(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	validReq := dto.GetCurrentUserRequest{
		UserID: "user_id_123",
	}

	t.Run("should get current user successfully", func(t *testing.T) {
		deps := setupGetCurrentUser(t)

		email, _ := domain.NewEmail("jane.doe@example.com")
		role, _ := domain.NewRole("CLIENT")
		user := domain.NewUser(
			"user_id_123",
			"Jane",
			"Doe",
			email,
			domain.NewPasswordHash("hashed_password"),
			role,
			now,
		)

		deps.userRepository.FindByIDFn = func(ctx context.Context, id string) (*domain.User, error) {
			if id != "user_id_123" {
				t.Errorf("expected user ID user_id_123, got %s", id)
			}
			return user, nil
		}

		resp, err := deps.sut.Execute(context.Background(), validReq)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if resp.ID != "user_id_123" {
			t.Errorf("expected ID user_id_123, got %s", resp.ID)
		}
		if resp.FirstName != "Jane" {
			t.Errorf("expected FirstName Jane, got %s", resp.FirstName)
		}
		if resp.LastName != "Doe" {
			t.Errorf("expected LastName Doe, got %s", resp.LastName)
		}
		if resp.Email != "jane.doe@example.com" {
			t.Errorf("expected Email jane.doe@example.com, got %s", resp.Email)
		}
		if resp.Role != "CLIENT" {
			t.Errorf("expected Role CLIENT, got %s", resp.Role)
		}
	})

	t.Run("should fail when user is not found", func(t *testing.T) {
		deps := setupGetCurrentUser(t)

		deps.userRepository.FindByIDFn = func(ctx context.Context, id string) (*domain.User, error) {
			return nil, domain.ErrUserNotFound
		}

		_, err := deps.sut.Execute(context.Background(), validReq)
		if err == nil {
			t.Error("expected error when user is not found, got nil")
		}
		if !errors.Is(err, domain.ErrUserNotFound) {
			t.Errorf("expected ErrUserNotFound, got %v", err)
		}
	})

	t.Run("should fail when user repository returns error", func(t *testing.T) {
		deps := setupGetCurrentUser(t)

		dbErr := errors.New("database connection error")
		deps.userRepository.FindByIDFn = func(ctx context.Context, id string) (*domain.User, error) {
			return nil, dbErr
		}

		_, err := deps.sut.Execute(context.Background(), validReq)
		if err == nil {
			t.Error("expected error from user repository, got nil")
		}
		if !errors.Is(err, dbErr) {
			t.Errorf("expected dbErr, got %v", err)
		}
	})
}
