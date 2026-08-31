package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultant/mocks"
)

type testGetConsultant struct {
	consultantRepo *mocks.MockConsultantRepository
	professionRepo *mocks.MockProfessionRepository
	sut            *GetConsultant
}

func setUpGetConsultant(t *testing.T) *testGetConsultant {
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
		FindByIDFn: func(ctx context.Context, id string) (*domain.Consultant, error) {
			return consultant, nil
		},
		FindByUserIDFn: func(ctx context.Context, userID string) (*domain.Consultant, error) {
			return consultant, nil
		},
	}

	profession := domain.NewProfession("prof_9ee432d7-b672-40ae-b03f-c1f1fb696621", "SOFTWARE_ENGINEER", now)
	professionRepo := &mocks.MockProfessionRepository{
		GetProfessionByIDFn: func(ctx context.Context, professionID string) (*domain.Profession, error) {
			return &profession, nil
		},
	}

	expertiseRepo := &mocks.MockExpertiseRepository{}

	sut := NewGetConsultantUsecase(consultantRepo, professionRepo, expertiseRepo)

	return &testGetConsultant{
		consultantRepo: consultantRepo,
		professionRepo: professionRepo,
		sut:            sut,
	}
}

func TestGetConsultant_ByID(t *testing.T) {
	t.Run("should get consultant by ID successfully", func(t *testing.T) {
		tc := setUpGetConsultant(t)

		res, err := tc.sut.ByID(context.Background(), "con_123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if res == nil {
			t.Fatal("expected response, got nil")
		}
		if res.ID != "con_123" {
			t.Errorf("expected ID con_123, got %s", res.ID)
		}
		if res.Profession != "SOFTWARE_ENGINEER" {
			t.Errorf("expected Profession SOFTWARE_ENGINEER, got %s", res.Profession)
		}
		if res.DisplayName != "Jane Doe Tech" {
			t.Errorf("expected DisplayName Jane Doe Tech, got %s", res.DisplayName)
		}
		if res.UserID != "user_123" {
			t.Errorf("expected UserID user_123, got %s", res.UserID)
		}
	})

	t.Run("should fail when consultant repository FindByID returns error", func(t *testing.T) {
		tc := setUpGetConsultant(t)
		expectedErr := errors.New("db error")
		tc.consultantRepo.FindByIDFn = func(ctx context.Context, id string) (*domain.Consultant, error) {
			return nil, expectedErr
		}

		res, err := tc.sut.ByID(context.Background(), "con_123")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
		if res != nil {
			t.Errorf("expected nil response, got %v", res)
		}
	})

	t.Run("should fail when profession repository GetProfessionByID returns error", func(t *testing.T) {
		tc := setUpGetConsultant(t)
		expectedErr := domain.ErrInvalidProfession
		tc.professionRepo.GetProfessionByIDFn = func(ctx context.Context, professionID string) (*domain.Profession, error) {
			return nil, expectedErr
		}

		res, err := tc.sut.ByID(context.Background(), "con_123")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
		if res != nil {
			t.Errorf("expected nil response, got %v", res)
		}
	})
}

func TestGetConsultant_ByUserID(t *testing.T) {
	t.Run("should get consultant by UserID successfully", func(t *testing.T) {
		tc := setUpGetConsultant(t)

		res, err := tc.sut.ByUserID(context.Background(), "user_123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if res == nil {
			t.Fatal("expected response, got nil")
		}
		if res.ID != "con_123" {
			t.Errorf("expected ID con_123, got %s", res.ID)
		}
		if res.Profession != "SOFTWARE_ENGINEER" {
			t.Errorf("expected Profession SOFTWARE_ENGINEER, got %s", res.Profession)
		}
	})

	t.Run("should fail when consultant repository FindByUserID returns error", func(t *testing.T) {
		tc := setUpGetConsultant(t)
		expectedErr := errors.New("db error")
		tc.consultantRepo.FindByUserIDFn = func(ctx context.Context, userID string) (*domain.Consultant, error) {
			return nil, expectedErr
		}

		res, err := tc.sut.ByUserID(context.Background(), "user_123")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
		if res != nil {
			t.Errorf("expected nil response, got %v", res)
		}
	})

	t.Run("should fail when profession repository GetProfessionByID returns error for ByUserID", func(t *testing.T) {
		tc := setUpGetConsultant(t)
		expectedErr := domain.ErrInvalidProfession
		tc.professionRepo.GetProfessionByIDFn = func(ctx context.Context, professionID string) (*domain.Profession, error) {
			return nil, expectedErr
		}

		res, err := tc.sut.ByUserID(context.Background(), "user_123")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
		if res != nil {
			t.Errorf("expected nil response, got %v", res)
		}
	})
}
