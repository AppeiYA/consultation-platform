package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultant/mocks"
	"github.com/AppeiYA/consultation-platform/internal/consultant/usecase/dto"
	shared_mocks "github.com/AppeiYA/consultation-platform/internal/shared/mocks"
)

type testRegisterConsultant struct {
	consultantRepo *mocks.MockConsultantRepository
	professionRepo *mocks.MockProfessionRepository
	roleAssigner   *mocks.MockRoleAssigner
	idGenerator    *shared_mocks.MockIDGenerator
	clock          *shared_mocks.MockClock

	sut *RegisterConsultant
}

func setUpRegisterConsultant(t *testing.T) *testRegisterConsultant {
	t.Helper()

	consultantRepo := &mocks.MockConsultantRepository{
		SaveFn: func(ctx context.Context, consultant *domain.Consultant) error {
			return nil
		},
		ExistsByUserIDFn: func(ctx context.Context, userID string) (bool, error) {
			return false, nil
		},
	}
	idGenerator := &shared_mocks.MockIDGenerator{
		GenerateFn: func(prefix string) (string, error) {
			return "con_123456", nil
		},
	}
	clock := &shared_mocks.MockClock{
		NowFn: func() time.Time {
			return time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
		},
	}
	professionRepo := &mocks.MockProfessionRepository{
		GetProfessionByIDFn: func(ctx context.Context, professionID string) (*domain.Profession, error) {
			prof := domain.NewProfession("prof_001", "SOFTWARE_ENGINEER", time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC))
			return &prof, nil
		},
	}
	roleAssigner := &mocks.MockRoleAssigner{
		AssignConsultantRoleFn: func(ctx context.Context, userID string) error {
			return nil
		},
	}

	expertiseRepo := &mocks.MockExpertiseRepository{}

	sut := NewRegisterConsultantUsecase(consultantRepo, professionRepo, expertiseRepo, roleAssigner, idGenerator, clock)

	return &testRegisterConsultant{
		consultantRepo: consultantRepo,
		professionRepo: professionRepo,
		roleAssigner:   roleAssigner,
		idGenerator:    idGenerator,
		clock:          clock,
		sut:            sut,
	}
}

func TestRegisterConsultant_Execute(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	userID := "user_123"
	validReq := &dto.RegisterConsultantDTO{
		ProfessionID:    "prof_001",
		DisplayName:     "John Doe",
		Bio:             "Experienced software engineer with 5 years in tech.",
		YearsExperience: 5,
	}

	t.Run("should register consultant successfully", func(t *testing.T) {
		tc := setUpRegisterConsultant(t)

		var generatedPrefix string
		tc.idGenerator.GenerateFn = func(prefix string) (string, error) {
			generatedPrefix = prefix
			return "con_123456", nil
		}

		var savedConsultant *domain.Consultant
		tc.consultantRepo.SaveFn = func(ctx context.Context, consultant *domain.Consultant) error {
			savedConsultant = consultant
			return nil
		}

		var assignedUserID string
		tc.roleAssigner.AssignConsultantRoleFn = func(ctx context.Context, uID string) error {
			assignedUserID = uID
			return nil
		}

		err := tc.sut.Execute(context.Background(), userID, validReq)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if generatedPrefix != domain.ConsultantIDPrefix {
			t.Errorf("expected ID prefix %s, got %s", domain.ConsultantIDPrefix, generatedPrefix)
		}

		if savedConsultant == nil {
			t.Fatal("expected saved consultant, got nil")
		}

		if savedConsultant.ID() != "con_123456" {
			t.Errorf("expected ID con_123456, got %s", savedConsultant.ID())
		}
		if savedConsultant.UserID() != userID {
			t.Errorf("expected UserID %s, got %s", userID, savedConsultant.UserID())
		}
		if !savedConsultant.IsAcceptingClients() {
			t.Errorf("expected IsAcceptingClients to be true")
		}
		if !savedConsultant.CreatedAt().Equal(now) {
			t.Errorf("expected CreatedAt %v, got %v", now, savedConsultant.CreatedAt())
		}
		if !savedConsultant.UpdatedAt().Equal(now) {
			t.Errorf("expected UpdatedAt %v, got %v", now, savedConsultant.UpdatedAt())
		}
		if assignedUserID != userID {
			t.Errorf("expected assigned user ID %s, got %s", userID, assignedUserID)
		}
	})

	t.Run("should fail when ExistsByUserID returns error", func(t *testing.T) {
		tc := setUpRegisterConsultant(t)
		expectedErr := errors.New("database connection failed")
		tc.consultantRepo.ExistsByUserIDFn = func(ctx context.Context, userID string) (bool, error) {
			return false, expectedErr
		}

		err := tc.sut.Execute(context.Background(), userID, validReq)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("should fail when consultant already exists", func(t *testing.T) {
		tc := setUpRegisterConsultant(t)
		tc.consultantRepo.ExistsByUserIDFn = func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		}

		err := tc.sut.Execute(context.Background(), userID, validReq)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrConsultantAlreadyExists) {
			t.Errorf("expected ErrConsultantAlreadyExists, got %v", err)
		}
	})

	t.Run("should fail when ID generator returns error", func(t *testing.T) {
		tc := setUpRegisterConsultant(t)
		expectedErr := errors.New("failed to generate ID")
		tc.idGenerator.GenerateFn = func(prefix string) (string, error) {
			return "", expectedErr
		}

		err := tc.sut.Execute(context.Background(), userID, validReq)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("should fail when profession is invalid", func(t *testing.T) {
		tc := setUpRegisterConsultant(t)
		req := *validReq
		req.ProfessionID = ""

		err := tc.sut.Execute(context.Background(), userID, &req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrInvalidProfession) {
			t.Errorf("expected ErrInvalidProfession, got %v", err)
		}
	})

	t.Run("should fail when profession does not exist", func(t *testing.T) {
		tc := setUpRegisterConsultant(t)
		expectedErr := domain.ErrInvalidProfession
		tc.professionRepo.GetProfessionByIDFn = func(ctx context.Context, professionID string) (*domain.Profession, error) {
			return nil, expectedErr
		}

		err := tc.sut.Execute(context.Background(), userID, validReq)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("should fail when display name is invalid", func(t *testing.T) {
		tc := setUpRegisterConsultant(t)
		req := *validReq
		req.DisplayName = "Short"

		err := tc.sut.Execute(context.Background(), userID, &req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrInvalidDisplayNameLength) {
			t.Errorf("expected ErrInvalidDisplayNameLength, got %v", err)
		}
	})

	t.Run("should fail when bio is invalid", func(t *testing.T) {
		tc := setUpRegisterConsultant(t)
		req := *validReq
		req.Bio = ""

		err := tc.sut.Execute(context.Background(), userID, &req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrEmptyBio) {
			t.Errorf("expected ErrEmptyBio, got %v", err)
		}
	})

	t.Run("should fail when years of experience is invalid", func(t *testing.T) {
		tc := setUpRegisterConsultant(t)
		req := *validReq
		req.YearsExperience = 0

		err := tc.sut.Execute(context.Background(), userID, &req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrInvalidYearsExperience) {
			t.Errorf("expected ErrInvalidYearsExperience, got %v", err)
		}
	})

	t.Run("should fail when consultant repository Save returns error", func(t *testing.T) {
		tc := setUpRegisterConsultant(t)
		expectedErr := errors.New("failed to save consultant")
		tc.consultantRepo.SaveFn = func(ctx context.Context, consultant *domain.Consultant) error {
			return expectedErr
		}

		err := tc.sut.Execute(context.Background(), userID, validReq)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("should fail when role assigner returns error", func(t *testing.T) {
		tc := setUpRegisterConsultant(t)
		expectedErr := errors.New("failed to assign role")
		tc.roleAssigner.AssignConsultantRoleFn = func(ctx context.Context, uID string) error {
			return expectedErr
		}

		err := tc.sut.Execute(context.Background(), userID, validReq)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})
}