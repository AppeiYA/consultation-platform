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
	sut *GetCurrentUser

	userRepository     *mocks.MockUserRepository
	sessionStore       *mocks.MockSessionStore
	sessionTokenHasher *mocks.MockSessionTokenHasher
	clock              *mocks.MockClock
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

	sessionStore := &mocks.MockSessionStore{
		FindByTokenHashFn: func(ctx context.Context, tokenHash string) (*domain.Session, error) {
			return nil, domain.ErrSessionNotFound
		},
	}

	sessionTokenHasher := &mocks.MockSessionTokenHasher{
		HashFn: func(token string) (string, error) {
			return "hashed_token", nil
		},
	}

	clock := &mocks.MockClock{
		NowFn: func() time.Time {
			return time.Now()
		},
	}

	sut := NewGetCurrentUser(
		sessionStore,
		sessionTokenHasher,
		userRepo,
		clock,
	)

	return &testGetCurrentUser{
		sut:                sut,
		userRepository:     userRepo,
		sessionStore:       sessionStore,
		sessionTokenHasher: sessionTokenHasher,
		clock:              clock,
	}
}

func TestGetCurrentUser_Execute(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	validSessionToken, _ := domain.NewSessionToken("12341234123412341234123412341234")
	validReq := dto.GetCurrentUserRequest{
		SessionToken: validSessionToken.String(),
	}

	t.Run("should get current user successfully", func(t *testing.T) {
		deps := setupGetCurrentUser(t)
		deps.clock.NowFn = func() time.Time { return now }

		session, err := domain.NewSession("session_id", "user_id_123", "hashed_token", now, 24*time.Hour)
		if err != nil {
			t.Fatalf("failed to create test session: %v", err)
		}

		deps.sessionStore.FindByTokenHashFn = func(ctx context.Context, tokenHash string) (*domain.Session, error) {
			if tokenHash != "hashed_token" {
				t.Errorf("expected tokenHash hashed_token, got %s", tokenHash)
			}
			return session, nil
		}

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

	t.Run("should fail on invalid session token format", func(t *testing.T) {
		deps := setupGetCurrentUser(t)

		invalidReq := dto.GetCurrentUserRequest{
			SessionToken: "invalid_token",
		}

		_, err := deps.sut.Execute(context.Background(), invalidReq)
		if err == nil {
			t.Error("expected error for invalid session token format, got nil")
		}
	})

	t.Run("should fail when session token hashing fails", func(t *testing.T) {
		deps := setupGetCurrentUser(t)

		hashErr := errors.New("hashing failed")
		deps.sessionTokenHasher.HashFn = func(token string) (string, error) {
			return "", hashErr
		}

		_, err := deps.sut.Execute(context.Background(), validReq)
		if err == nil {
			t.Error("expected error when hashing fails, got nil")
		}
		if !errors.Is(err, hashErr) {
			t.Errorf("expected hashErr, got %v", err)
		}
	})

	t.Run("should fail when session is not found", func(t *testing.T) {
		deps := setupGetCurrentUser(t)

		deps.sessionStore.FindByTokenHashFn = func(ctx context.Context, tokenHash string) (*domain.Session, error) {
			return nil, domain.ErrSessionNotFound
		}

		_, err := deps.sut.Execute(context.Background(), validReq)
		if err == nil {
			t.Error("expected error when session is not found, got nil")
		}
		if !errors.Is(err, domain.ErrSessionNotFound) {
			t.Errorf("expected ErrSessionNotFound, got %v", err)
		}
	})

	t.Run("should fail when session is expired", func(t *testing.T) {
		deps := setupGetCurrentUser(t)
		deps.clock.NowFn = func() time.Time { return now }

		createdAt := now.Add(-2 * time.Hour)
		expiredSession, err := domain.NewSession("session_id", "user_id_123", "hashed_token", createdAt, 1*time.Hour)
		if err != nil {
			t.Fatalf("failed to create expired session: %v", err)
		}

		deps.sessionStore.FindByTokenHashFn = func(ctx context.Context, tokenHash string) (*domain.Session, error) {
			return expiredSession, nil
		}

		_, err = deps.sut.Execute(context.Background(), validReq)
		if err == nil {
			t.Error("expected error for expired session, got nil")
		}
		if !errors.Is(err, domain.ErrSessionExpired) {
			t.Errorf("expected ErrSessionExpired, got %v", err)
		}
	})

	t.Run("should fail when user is not found", func(t *testing.T) {
		deps := setupGetCurrentUser(t)
		deps.clock.NowFn = func() time.Time { return now }

		session, _ := domain.NewSession("session_id", "user_id_123", "hashed_token", now, 24*time.Hour)
		deps.sessionStore.FindByTokenHashFn = func(ctx context.Context, tokenHash string) (*domain.Session, error) {
			return session, nil
		}

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
		deps.clock.NowFn = func() time.Time { return now }

		session, _ := domain.NewSession("session_id", "user_id_123", "hashed_token", now, 24*time.Hour)
		deps.sessionStore.FindByTokenHashFn = func(ctx context.Context, tokenHash string) (*domain.Session, error) {
			return session, nil
		}

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