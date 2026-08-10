package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/identity/domain"
	shared_mocks "github.com/AppeiYA/consultation-platform/internal/shared/mocks"
	"github.com/AppeiYA/consultation-platform/internal/identity/mocks"
	"github.com/AppeiYA/consultation-platform/internal/identity/usecase/dto"
)

type testLoginUser struct {
	userRepo              *mocks.MockUserRepository
	sessionStore          *mocks.MockSessionStore
	passwordHasher        *mocks.MockPasswordHasher
	sessionTokenHasher    *mocks.MockSessionTokenHasher
	sessionTokenGenerator *mocks.MockSessionTokenGenerator
	idGenerator           *shared_mocks.MockIDGenerator
	clock                 *shared_mocks.MockClock 
	t                     *testing.T
}

func setupLoginUser(t *testing.T) testLoginUser {
	return testLoginUser{
		userRepo: &mocks.MockUserRepository{
			FindByEmailFn: func(ctx context.Context, email domain.Email) (*domain.User, error) {
				return nil, domain.ErrUserNotFound
			},
			SaveFn: func(ctx context.Context, user *domain.User) error {
				return nil
			},
		},
		sessionStore: &mocks.MockSessionStore{
			SaveFn: func(ctx context.Context, session *domain.Session) error {
				return nil
			},
			FindByTokenHashFn: func(ctx context.Context, tokenHash string) (*domain.Session, error) {
				return nil, domain.ErrSessionNotFound
			},
			DeleteFn: func(ctx context.Context, sessionID string) error {
				return nil
			},
		},
		sessionTokenHasher: &mocks.MockSessionTokenHasher{
			HashFn: func(token string) (string, error) {
				return "", nil
			},
		},
		sessionTokenGenerator: &mocks.MockSessionTokenGenerator{
			GenerateFn: func() (domain.SessionToken, error) {
				return domain.SessionToken{}, nil
			},
		},
		passwordHasher: &mocks.MockPasswordHasher{
			HashFn: func(password string) (string, error) {
				return "hashed_password", nil
			},
			CompareFn: func(
				password string,
				hash string,
			) (bool, error) {
				return hash == password, nil
			},
		},
		idGenerator: &shared_mocks.MockIDGenerator{
			GenerateFn: func(_ string) (string, error) {
				return "user_id", nil
			},
		},
		clock: &shared_mocks.MockClock{
			NowFn: func() time.Time {
				return time.Now()
			},
		},
		t: t,
	}
}

func TestLoginUser_Execute(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	sessionTTL := 24 * time.Hour
	validParams := dto.LoginRequest{
		Email:    "appeimisic@gmail.com",
		Password: "testLogin123$",
	}

	t.Run("should login user successfully", func(t *testing.T) {
		deps := setupLoginUser(t)

		deps.clock.NowFn = func() time.Time {
			return now
		}

		deps.userRepo.FindByEmailFn = func(
			ctx context.Context,
			email domain.Email,
		) (*domain.User, error) {
			role, _ := domain.NewRole("CLIENT")
			return domain.NewUser(
				"user_id",
				"John",
				"Doe",
				email,
				domain.NewPasswordHash("hashed_password"),
				role,
				now,
			), nil
		}

		deps.passwordHasher.CompareFn = func(
			plain string,
			hashed string,
		) (bool, error) {
			return true, nil
		}

		expectedToken, err := domain.NewSessionToken("12345678901234567890123456789012")
		if err != nil {
			t.Fatalf("failed to create session token: %v", err)
		}

		deps.sessionTokenGenerator.GenerateFn = func() (domain.SessionToken, error) {
			return expectedToken, nil
		}

		deps.sessionTokenHasher.HashFn = func(token string) (string, error) {
			return "hashed_session_token", nil
		}

		uc := NewLoginUser(
			deps.userRepo,
			deps.sessionStore,
			deps.passwordHasher,
			deps.sessionTokenHasher,
			deps.sessionTokenGenerator,
			deps.idGenerator,
			deps.clock,
			sessionTTL,
		)

		result, err := uc.Execute(context.Background(), validParams)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result.SessionToken != expectedToken {
			t.Errorf("expected session token %v, got %v", expectedToken, result.SessionToken)
		}
	})

	t.Run("should fail on invalid email", func(t *testing.T) {
		deps := setupLoginUser(t)

		invalidEmailReq := validParams
		invalidEmailReq.Email = "invalid-email"

		uc := NewLoginUser(
			deps.userRepo,
			deps.sessionStore,
			deps.passwordHasher,
			deps.sessionTokenHasher,
			deps.sessionTokenGenerator,
			deps.idGenerator,
			deps.clock,
			sessionTTL,
		)

		_, err := uc.Execute(context.Background(), invalidEmailReq)
		if err == nil {
			t.Error("expected error for invalid email, got nil")
		}
	})

	t.Run("should fail on invalid password format", func(t *testing.T) {
		deps := setupLoginUser(t)
		invalidPasswordReq := validParams
		invalidPasswordReq.Password = "weak"

		uc := NewLoginUser(
			deps.userRepo,
			deps.sessionStore,
			deps.passwordHasher,
			deps.sessionTokenHasher,
			deps.sessionTokenGenerator,
			deps.idGenerator,
			deps.clock,
			sessionTTL,
		)

		_, err := uc.Execute(context.Background(), invalidPasswordReq)
		if err == nil {
			t.Error("expected error for invalid password format, got nil")
		}
	})

	t.Run("should fail when user is not found", func(t *testing.T) {
		deps := setupLoginUser(t)
		deps.userRepo.FindByEmailFn = func(ctx context.Context, email domain.Email) (*domain.User, error) {
			return nil, domain.ErrUserNotFound
		}

		uc := NewLoginUser(
			deps.userRepo,
			deps.sessionStore,
			deps.passwordHasher,
			deps.sessionTokenHasher,
			deps.sessionTokenGenerator,
			deps.idGenerator,
			deps.clock,
			sessionTTL,
		)

		_, err := uc.Execute(context.Background(), validParams)
		if err == nil {
			t.Error("expected error when user is not found, got nil")
		}
		if !errors.Is(err, domain.ErrUserNotFound) {
			t.Errorf("expected ErrUserNotFound, got %v", err)
		}
	})

	t.Run("should fail when password compare fails with error", func(t *testing.T) {
		deps := setupLoginUser(t)
		deps.userRepo.FindByEmailFn = func(ctx context.Context, email domain.Email) (*domain.User, error) {
			role, _ := domain.NewRole("CLIENT")
			return domain.NewUser("user_id", "John", "Doe", email, domain.NewPasswordHash("hashed_password"), role, now), nil
		}
		compareErr := errors.New("compare failed")
		deps.passwordHasher.CompareFn = func(password, hash string) (bool, error) {
			return false, compareErr
		}

		uc := NewLoginUser(
			deps.userRepo,
			deps.sessionStore,
			deps.passwordHasher,
			deps.sessionTokenHasher,
			deps.sessionTokenGenerator,
			deps.idGenerator,
			deps.clock,
			sessionTTL,
		)

		_, err := uc.Execute(context.Background(), validParams)
		if err == nil {
			t.Error("expected error when password compare fails, got nil")
		}
		if !errors.Is(err, compareErr) {
			t.Errorf("expected compareErr, got %v", err)
		}
	})

	t.Run("should fail when password does not match", func(t *testing.T) {
		deps := setupLoginUser(t)
		deps.userRepo.FindByEmailFn = func(ctx context.Context, email domain.Email) (*domain.User, error) {
			role, _ := domain.NewRole("CLIENT")
			return domain.NewUser("user_id", "John", "Doe", email, domain.NewPasswordHash("hashed_password"), role, now), nil
		}
		deps.passwordHasher.CompareFn = func(password, hash string) (bool, error) {
			return false, nil
		}

		uc := NewLoginUser(
			deps.userRepo,
			deps.sessionStore,
			deps.passwordHasher,
			deps.sessionTokenHasher,
			deps.sessionTokenGenerator,
			deps.idGenerator,
			deps.clock,
			sessionTTL,
		)

		_, err := uc.Execute(context.Background(), validParams)
		if err == nil {
			t.Error("expected error when password does not match, got nil")
		}
		if !errors.Is(err, domain.ErrInvalidPassword) {
			t.Errorf("expected ErrInvalidPassword, got %v", err)
		}
	})

	t.Run("should fail when session token generation fails", func(t *testing.T) {
		deps := setupLoginUser(t)
		deps.userRepo.FindByEmailFn = func(ctx context.Context, email domain.Email) (*domain.User, error) {
			role, _ := domain.NewRole("CLIENT")
			return domain.NewUser("user_id", "John", "Doe", email, domain.NewPasswordHash("hashed_password"), role, now), nil
		}
		deps.passwordHasher.CompareFn = func(password, hash string) (bool, error) {
			return true, nil
		}
		genErr := errors.New("token generation failed")
		deps.sessionTokenGenerator.GenerateFn = func() (domain.SessionToken, error) {
			return domain.SessionToken{}, genErr
		}

		uc := NewLoginUser(
			deps.userRepo,
			deps.sessionStore,
			deps.passwordHasher,
			deps.sessionTokenHasher,
			deps.sessionTokenGenerator,
			deps.idGenerator,
			deps.clock,
			sessionTTL,
		)

		_, err := uc.Execute(context.Background(), validParams)
		if err == nil {
			t.Error("expected error when token generation fails, got nil")
		}
		if !errors.Is(err, genErr) {
			t.Errorf("expected genErr, got %v", err)
		}
	})

	t.Run("should fail when session token hashing fails", func(t *testing.T) {
		deps := setupLoginUser(t)
		deps.userRepo.FindByEmailFn = func(ctx context.Context, email domain.Email) (*domain.User, error) {
			role, _ := domain.NewRole("CLIENT")
			return domain.NewUser("user_id", "John", "Doe", email, domain.NewPasswordHash("hashed_password"), role, now), nil
		}
		deps.passwordHasher.CompareFn = func(password, hash string) (bool, error) {
			return true, nil
		}
		validToken, _ := domain.NewSessionToken("12345678901234567890123456789012")
		deps.sessionTokenGenerator.GenerateFn = func() (domain.SessionToken, error) {
			return validToken, nil
		}
		hashErr := errors.New("token hashing failed")
		deps.sessionTokenHasher.HashFn = func(token string) (string, error) {
			return "", hashErr
		}

		uc := NewLoginUser(
			deps.userRepo,
			deps.sessionStore,
			deps.passwordHasher,
			deps.sessionTokenHasher,
			deps.sessionTokenGenerator,
			deps.idGenerator,
			deps.clock,
			sessionTTL,
		)

		_, err := uc.Execute(context.Background(), validParams)
		if err == nil {
			t.Error("expected error when token hashing fails, got nil")
		}
		if !errors.Is(err, hashErr) {
			t.Errorf("expected hashErr, got %v", err)
		}
	})

	t.Run("should fail when session ID generation fails", func(t *testing.T) {
		deps := setupLoginUser(t)
		deps.userRepo.FindByEmailFn = func(ctx context.Context, email domain.Email) (*domain.User, error) {
			role, _ := domain.NewRole("CLIENT")
			return domain.NewUser("user_id", "John", "Doe", email, domain.NewPasswordHash("hashed_password"), role, now), nil
		}
		deps.passwordHasher.CompareFn = func(password, hash string) (bool, error) {
			return true, nil
		}
		validToken, _ := domain.NewSessionToken("12345678901234567890123456789012")
		deps.sessionTokenGenerator.GenerateFn = func() (domain.SessionToken, error) {
			return validToken, nil
		}
		deps.sessionTokenHasher.HashFn = func(token string) (string, error) {
			return "hashed_token", nil
		}
		idGenErr := errors.New("session ID generation failed")
		deps.idGenerator.GenerateFn = func(_ string) (string, error) {
			return "", idGenErr
		}

		uc := NewLoginUser(
			deps.userRepo,
			deps.sessionStore,
			deps.passwordHasher,
			deps.sessionTokenHasher,
			deps.sessionTokenGenerator,
			deps.idGenerator,
			deps.clock,
			sessionTTL,
		)

		_, err := uc.Execute(context.Background(), validParams)
		if err == nil {
			t.Error("expected error when ID generation fails, got nil")
		}
		if !errors.Is(err, idGenErr) {
			t.Errorf("expected idGenErr, got %v", err)
		}
	})

	t.Run("should fail when session store save fails", func(t *testing.T) {
		deps := setupLoginUser(t)
		deps.userRepo.FindByEmailFn = func(ctx context.Context, email domain.Email) (*domain.User, error) {
			role, _ := domain.NewRole("CLIENT")
			return domain.NewUser("user_id", "John", "Doe", email, domain.NewPasswordHash("hashed_password"), role, now), nil
		}
		deps.passwordHasher.CompareFn = func(password, hash string) (bool, error) {
			return true, nil
		}
		validToken, _ := domain.NewSessionToken("12345678901234567890123456789012")
		deps.sessionTokenGenerator.GenerateFn = func() (domain.SessionToken, error) {
			return validToken, nil
		}
		deps.sessionTokenHasher.HashFn = func(token string) (string, error) {
			return "hashed_token", nil
		}
		saveErr := errors.New("session save failed")
		deps.sessionStore.SaveFn = func(ctx context.Context, session *domain.Session) error {
			return saveErr
		}

		uc := NewLoginUser(
			deps.userRepo,
			deps.sessionStore,
			deps.passwordHasher,
			deps.sessionTokenHasher,
			deps.sessionTokenGenerator,
			deps.idGenerator,
			deps.clock,
			sessionTTL,
		)

		_, err := uc.Execute(context.Background(), validParams)
		if err == nil {
			t.Error("expected error when session save fails, got nil")
		}
		if !errors.Is(err, saveErr) {
			t.Errorf("expected saveErr, got %v", err)
		}
	})
}