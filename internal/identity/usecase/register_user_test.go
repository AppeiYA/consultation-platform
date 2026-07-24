package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/identity/domain"
	"github.com/AppeiYA/consultation-platform/internal/identity/mocks"
)

type TestRegisterUser struct {
	userRepo       *mocks.MockUserRepository
	passwordHasher *mocks.MockPasswordHasher
	idGenerator    *mocks.MockIDGenerator
	clock          *mocks.MockClock
	t              *testing.T
}

func setup(t *testing.T) TestRegisterUser {
	return TestRegisterUser{
		userRepo: &mocks.MockUserRepository{
			FindByEmailFn: func(ctx context.Context, email domain.Email) (*domain.User, error) {
				return nil, domain.ErrUserNotFound
			},
			SaveFn: func(ctx context.Context, user *domain.User) error {
				return nil
			},
		},
		passwordHasher: &mocks.MockPasswordHasher{
			HashFn: func(password string) (string, error) {
				return "hashed_password", nil
			},
		},
		idGenerator: &mocks.MockIDGenerator{
			GenerateFn: func() (string, error) {
				return "user_id", nil
			},
		},
		clock: &mocks.MockClock{
			NowFn: func() time.Time {
				return time.Now()
			},
		},
		t: t,
	}
}

func TestRegisterUser_Execute(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	validParams := RegisterUserParams{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "SecurePass123!",
		Role:      "CLIENT",
	}

	t.Run("should register user successfully", func(t *testing.T) {
		deps := setup(t)
		deps.clock.NowFn = func() time.Time { return now }

		var savedUser *domain.User
		deps.userRepo.SaveFn = func(ctx context.Context, user *domain.User) error {
			savedUser = user
			return nil
		}

		uc := NewRegisterUser(deps.userRepo, deps.passwordHasher, deps.idGenerator, deps.clock)
		resp, err := uc.Execute(context.Background(), validParams)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if resp.UserID != "user_id" {
			t.Errorf("expected UserID user_id, got %s", resp.UserID)
		}
		if savedUser == nil {
			t.Fatal("expected saved user, got nil")
		}
		if savedUser.ID() != "user_id" {
			t.Errorf("expected saved user ID user_id, got %s", savedUser.ID())
		}
		if savedUser.Email().String() != validParams.Email {
			t.Errorf("expected saved email %s, got %s", validParams.Email, savedUser.Email().String())
		}
	})

	t.Run("should fail when email is invalid", func(t *testing.T) {
		deps := setup(t)
		uc := NewRegisterUser(deps.userRepo, deps.passwordHasher, deps.idGenerator, deps.clock)

		params := validParams
		params.Email = "invalid-email"

		resp, err := uc.Execute(context.Background(), params)
		if err == nil {
			t.Error("expected error for invalid email, got nil")
		}
		if resp != nil {
			t.Errorf("expected nil response, got %v", resp)
		}
	})

	t.Run("should fail when password is invalid", func(t *testing.T) {
		deps := setup(t)
		uc := NewRegisterUser(deps.userRepo, deps.passwordHasher, deps.idGenerator, deps.clock)

		params := validParams
		params.Password = "weak"

		resp, err := uc.Execute(context.Background(), params)
		if err == nil {
			t.Error("expected error for invalid password, got nil")
		}
		if resp != nil {
			t.Errorf("expected nil response, got %v", resp)
		}
	})

	t.Run("should fail when role is invalid", func(t *testing.T) {
		deps := setup(t)
		uc := NewRegisterUser(deps.userRepo, deps.passwordHasher, deps.idGenerator, deps.clock)

		params := validParams
		params.Role = "INVALID_ROLE"

		resp, err := uc.Execute(context.Background(), params)
		if err == nil {
			t.Error("expected error for invalid role, got nil")
		}
		if resp != nil {
			t.Errorf("expected nil response, got %v", resp)
		}
	})

	t.Run("should fail when user already exists", func(t *testing.T) {
		deps := setup(t)
		deps.userRepo.FindByEmailFn = func(ctx context.Context, email domain.Email) (*domain.User, error) {
			existingEmail, _ := domain.NewEmail(validParams.Email)
			existingPasswordHash := domain.NewPasswordHash("hashed")
			existingRole, _ := domain.NewRole("CLIENT")
			return domain.NewUser("existing_id", "John", "Doe", existingEmail, existingPasswordHash, existingRole, now), nil
		}

		uc := NewRegisterUser(deps.userRepo, deps.passwordHasher, deps.idGenerator, deps.clock)
		resp, err := uc.Execute(context.Background(), validParams)

		if err == nil {
			t.Error("expected error when user already exists, got nil")
		}
		if !errors.Is(err, domain.ErrUserAlreadyExists) {
			t.Errorf("expected ErrUserAlreadyExists, got %v", err)
		}
		if resp != nil {
			t.Errorf("expected nil response, got %v", resp)
		}
	})

	t.Run("should fail when FindByEmail returns database error", func(t *testing.T) {
		deps := setup(t)
		dbErr := errors.New("database error")
		deps.userRepo.FindByEmailFn = func(ctx context.Context, email domain.Email) (*domain.User, error) {
			return nil, dbErr
		}

		uc := NewRegisterUser(deps.userRepo, deps.passwordHasher, deps.idGenerator, deps.clock)
		resp, err := uc.Execute(context.Background(), validParams)

		if err == nil {
			t.Error("expected database error, got nil")
		}
		if !errors.Is(err, dbErr) {
			t.Errorf("expected dbErr, got %v", err)
		}
		if resp != nil {
			t.Errorf("expected nil response, got %v", resp)
		}
	})

	t.Run("should fail when password hashing fails", func(t *testing.T) {
		deps := setup(t)
		hashErr := errors.New("hashing failed")
		deps.passwordHasher.HashFn = func(password string) (string, error) {
			return "", hashErr
		}

		uc := NewRegisterUser(deps.userRepo, deps.passwordHasher, deps.idGenerator, deps.clock)
		resp, err := uc.Execute(context.Background(), validParams)

		if err == nil {
			t.Error("expected hashing error, got nil")
		}
		if !errors.Is(err, hashErr) {
			t.Errorf("expected hashErr, got %v", err)
		}
		if resp != nil {
			t.Errorf("expected nil response, got %v", resp)
		}
	})

	t.Run("should fail when id generation fails", func(t *testing.T) {
		deps := setup(t)
		genErr := errors.New("id generation failed")
		deps.idGenerator.GenerateFn = func() (string, error) {
			return "", genErr
		}

		uc := NewRegisterUser(deps.userRepo, deps.passwordHasher, deps.idGenerator, deps.clock)
		resp, err := uc.Execute(context.Background(), validParams)

		if err == nil {
			t.Error("expected id generation error, got nil")
		}
		if !errors.Is(err, genErr) {
			t.Errorf("expected genErr, got %v", err)
		}
		if resp != nil {
			t.Errorf("expected nil response, got %v", resp)
		}
	})

	t.Run("should fail when userRepo save fails", func(t *testing.T) {
		deps := setup(t)
		saveErr := errors.New("save failed")
		deps.userRepo.SaveFn = func(ctx context.Context, user *domain.User) error {
			return saveErr
		}

		uc := NewRegisterUser(deps.userRepo, deps.passwordHasher, deps.idGenerator, deps.clock)
		resp, err := uc.Execute(context.Background(), validParams)

		if err == nil {
			t.Error("expected save error, got nil")
		}
		if !errors.Is(err, saveErr) {
			t.Errorf("expected saveErr, got %v", err)
		}
		if resp != nil {
			t.Errorf("expected nil response, got %v", resp)
		}
	})
}


