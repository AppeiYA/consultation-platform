package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/AppeiYA/consultation-platform/internal/identity/domain"
	"github.com/AppeiYA/consultation-platform/internal/identity/mocks"
	"github.com/AppeiYA/consultation-platform/internal/identity/usecase/dto"
)

type testLogoutUser struct {
	sut *LogoutUser

	sessionStore       *mocks.MockSessionStore
	sessionTokenHasher *mocks.MockSessionTokenHasher
}

func setupTestLogout(t *testing.T) *testLogoutUser {
	t.Helper()

	sessionStore := &mocks.MockSessionStore{
		DeleteFn: func(ctx context.Context, sessionID string) error {
			return nil
		},
	}

	sessionTokenHasher := &mocks.MockSessionTokenHasher{
		HashFn: func(token string) (string, error) {
			return "hashed_token", nil
		},
	}

	sut := NewLogoutUser(sessionStore, sessionTokenHasher)

	return &testLogoutUser{
		sut:                sut,
		sessionStore:       sessionStore,
		sessionTokenHasher: sessionTokenHasher,
	}
}

func TestLogout_Execute(t *testing.T) {
	validSessionToken, _ := domain.NewSessionToken("12341234123412341234123412341234")
	validReq := dto.LogoutRequest{
		SessionToken: validSessionToken.String(),
	}

	t.Run("should logout user successfully", func(t *testing.T) {
		deps := setupTestLogout(t)

		var deletedHash string
		deps.sessionStore.DeleteFn = func(ctx context.Context, sessionID string) error {
			deletedHash = sessionID
			return nil
		}

		var hashedToken string
		deps.sessionTokenHasher.HashFn = func(token string) (string, error) {
			hashedToken = token
			return "hashed_token_value", nil
		}

		_, err := deps.sut.Execute(context.Background(), validReq)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if hashedToken != validSessionToken.String() {
			t.Errorf("expected hashed token %s, got %s", validSessionToken.String(), hashedToken)
		}

		if deletedHash != "hashed_token_value" {
			t.Errorf("expected deleted session hash %s, got %s", "hashed_token_value", deletedHash)
		}
	})

	t.Run("should fail on invalid session token format", func(t *testing.T) {
		deps := setupTestLogout(t)

		invalidReq := dto.LogoutRequest{
			SessionToken: "invalid_token",
		}

		_, err := deps.sut.Execute(context.Background(), invalidReq)
		if err == nil {
			t.Error("expected error for invalid session token format, got nil")
		}
	})

	t.Run("should fail when session token hashing fails", func(t *testing.T) {
		deps := setupTestLogout(t)

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

	t.Run("should fail when session store delete fails", func(t *testing.T) {
		deps := setupTestLogout(t)

		deleteErr := errors.New("delete failed")
		deps.sessionStore.DeleteFn = func(ctx context.Context, sessionID string) error {
			return deleteErr
		}

		_, err := deps.sut.Execute(context.Background(), validReq)
		if err == nil {
			t.Error("expected error when session delete fails, got nil")
		}
		if !errors.Is(err, deleteErr) {
			t.Errorf("expected deleteErr, got %v", err)
		}
	})
}


