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

type validateSessionTest struct {
	sut *ValidateSession

	sessionStore       *mocks.MockSessionStore
	sessionTokenHasher *mocks.MockSessionTokenHasher
	clock              *mocks.MockClock
}

func setupValidateSessionTest(t *testing.T) *validateSessionTest {
	t.Helper()

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

	sut := NewValidateSession(
		sessionStore,
		sessionTokenHasher,
		clock,
	)
	return &validateSessionTest{
		sut:                sut,
		sessionStore:       sessionStore,
		sessionTokenHasher: sessionTokenHasher,
		clock:              clock,
	}
}

func TestValidateSession_Execute(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	validSessionToken, _ := domain.NewSessionToken("12341234123412341234123412341234")
	validReq := dto.ValidateSessionRequest{
		SessionToken: validSessionToken.String(),
	}

	t.Run("should validate session successfully", func(t *testing.T) {
		deps := setupValidateSessionTest(t)
		deps.clock.NowFn = func() time.Time { return now }

		session, err := domain.NewSession("session_id", "user_id_123", "jane.doe@example.com", "CLIENT", "hashed_token", now, 24*time.Hour)
		if err != nil {
			t.Fatalf("failed to create test session: %v", err)
		}

		deps.sessionStore.FindByTokenHashFn = func(ctx context.Context, tokenHash string) (*domain.Session, error) {
			if tokenHash != "hashed_token" {
				t.Errorf("expected tokenHash hashed_token, got %s", tokenHash)
			}
			return session, nil
		}

		resp, err := deps.sut.Execute(context.Background(), validReq)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if resp.SessionClaims.UserID != "user_id_123" {
			t.Errorf("expected UserID user_id_123, got %s", resp.SessionClaims.UserID)
		}
		if resp.SessionClaims.Email != "jane.doe@example.com" {
			t.Errorf("expected Email jane.doe@example.com, got %s", resp.SessionClaims.Email)
		}
		if resp.SessionClaims.Role != "CLIENT" {
			t.Errorf("expected Role CLIENT, got %s", resp.SessionClaims.Role)
		}
	})

	t.Run("should fail on invalid session token format", func(t *testing.T) {
		deps := setupValidateSessionTest(t)

		invalidReq := dto.ValidateSessionRequest{
			SessionToken: "invalid_token",
		}

		_, err := deps.sut.Execute(context.Background(), invalidReq)
		if err == nil {
			t.Error("expected error for invalid session token format, got nil")
		}
	})

	t.Run("should fail when session token hashing fails", func(t *testing.T) {
		deps := setupValidateSessionTest(t)

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
		deps := setupValidateSessionTest(t)

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
		deps := setupValidateSessionTest(t)
		deps.clock.NowFn = func() time.Time { return now }

		createdAt := now.Add(-2 * time.Hour)
		expiredSession, err := domain.NewSession("session_id", "user_id_123", "jane.doe@example.com", "CLIENT", "hashed_token", createdAt, 1*time.Hour)
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
}