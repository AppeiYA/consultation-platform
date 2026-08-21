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

type testUpdateConsultant struct {
	consultantRepo *mocks.MockConsultantRepository
	professionRepo *mocks.MockProfessionRepository
	clock          *shared_mocks.MockClock
	sut            *UpdateConsultant
}

func setUpUpdateConsultant(t *testing.T) *testUpdateConsultant {
	t.Helper()

	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	consultant, _ := domain.ReconstitueConsultant(
		"con_123",
		"user_123",
		"prof_9ee432d7-b672-40ae-b03f-c1f1fb696621",
		"Jane Doe Tech",
		"Experienced software engineer.",
		10,
		true,
		now,
		now,
	)

	consultantRepo := &mocks.MockConsultantRepository{
		FindByUserIDFn: func(ctx context.Context, userID string) (*domain.Consultant, error) {
			return consultant, nil
		},
		UpdateFn: func(ctx context.Context, consultant *domain.Consultant) error {
			return nil
		},
	}

	profession := domain.NewProfession("prof_12d965f5-e1f5-49aa-ac57-856772d236ce", "LAWYER", now)
	professionRepo := &mocks.MockProfessionRepository{
		GetProfessionByIDFn: func(ctx context.Context, professionID string) (*domain.Profession, error) {
			return &profession, nil
		},
	}

	clock := &shared_mocks.MockClock{
		NowFn: func() time.Time {
			return time.Date(2026, 1, 2, 12, 0, 0, 0, time.UTC)
		},
	}

	sut := NewUpdateConsultantUsecase(consultantRepo, professionRepo, clock)

	return &testUpdateConsultant{
		consultantRepo: consultantRepo,
		professionRepo: professionRepo,
		clock:          clock,
		sut:            sut,
	}
}

func TestUpdateConsultant_Execute(t *testing.T) {
	validReq := dto.UpdateConsultantDTO{
		ProfessionID:    "prof_12d965f5-e1f5-49aa-ac57-856772d236ce",
		DisplayName:     "Jane Doe Legal",
		Bio:             "Senior corporate lawyer with over 12 years experience.",
		YearsExperience: 12,
	}

	t.Run("should update consultant profile successfully", func(t *testing.T) {
		tc := setUpUpdateConsultant(t)

		var updatedConsultant *domain.Consultant
		tc.consultantRepo.UpdateFn = func(ctx context.Context, consultant *domain.Consultant) error {
			updatedConsultant = consultant
			return nil
		}

		err := tc.sut.Execute(context.Background(), "user_123", validReq)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if updatedConsultant == nil {
			t.Fatal("expected updated consultant, got nil")
		}
		if updatedConsultant.ProfessionID() != "prof_12d965f5-e1f5-49aa-ac57-856772d236ce" {
			t.Errorf("expected ProfessionID prof_12d965f5-e1f5-49aa-ac57-856772d236ce, got %s", updatedConsultant.ProfessionID())
		}
		if updatedConsultant.DisplayName().String() != "Jane Doe Legal" {
			t.Errorf("expected DisplayName Jane Doe Legal, got %s", updatedConsultant.DisplayName().String())
		}
		if updatedConsultant.YearsExperience().Int() != 12 {
			t.Errorf("expected YearsExperience 12, got %d", updatedConsultant.YearsExperience().Int())
		}
	})

	t.Run("should fail when consultant repository FindByUserID returns error", func(t *testing.T) {
		tc := setUpUpdateConsultant(t)
		expectedErr := domain.ErrConsultantNotFound
		tc.consultantRepo.FindByUserIDFn = func(ctx context.Context, userID string) (*domain.Consultant, error) {
			return nil, expectedErr
		}

		err := tc.sut.Execute(context.Background(), "user_123", validReq)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("should fail when profession is empty", func(t *testing.T) {
		tc := setUpUpdateConsultant(t)
		req := validReq
		req.ProfessionID = ""

		err := tc.sut.Execute(context.Background(), "user_123", req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrInvalidProfession) {
			t.Errorf("expected ErrInvalidProfession, got %v", err)
		}
	})

	t.Run("should fail when profession does not exist", func(t *testing.T) {
		tc := setUpUpdateConsultant(t)
		expectedErr := domain.ErrInvalidProfession
		tc.professionRepo.GetProfessionByIDFn = func(ctx context.Context, professionID string) (*domain.Profession, error) {
			return nil, expectedErr
		}

		err := tc.sut.Execute(context.Background(), "user_123", validReq)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("should fail when display name is invalid", func(t *testing.T) {
		tc := setUpUpdateConsultant(t)
		req := validReq
		req.DisplayName = "Short"

		err := tc.sut.Execute(context.Background(), "user_123", req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrInvalidDisplayNameLength) {
			t.Errorf("expected ErrInvalidDisplayNameLength, got %v", err)
		}
	})

	t.Run("should fail when bio is empty", func(t *testing.T) {
		tc := setUpUpdateConsultant(t)
		req := validReq
		req.Bio = ""

		err := tc.sut.Execute(context.Background(), "user_123", req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrEmptyBio) {
			t.Errorf("expected ErrEmptyBio, got %v", err)
		}
	})

	t.Run("should fail when years of experience is invalid", func(t *testing.T) {
		tc := setUpUpdateConsultant(t)
		req := validReq
		req.YearsExperience = 0

		err := tc.sut.Execute(context.Background(), "user_123", req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrInvalidYearsExperience) {
			t.Errorf("expected ErrInvalidYearsExperience, got %v", err)
		}
	})

	t.Run("should fail when consultant repository Update returns error", func(t *testing.T) {
		tc := setUpUpdateConsultant(t)
		expectedErr := errors.New("update failed")
		tc.consultantRepo.UpdateFn = func(ctx context.Context, consultant *domain.Consultant) error {
			return expectedErr
		}

		err := tc.sut.Execute(context.Background(), "user_123", validReq)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})
}
