package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultant/mocks"
	shared_mocks "github.com/AppeiYA/consultation-platform/internal/shared/mocks"
)

type testDeactivateAvailability struct {
	consultantRepo   *mocks.MockConsultantRepository
	availabilityRepo *mocks.MockAvailabilityRepository
	clock            *shared_mocks.MockClock

	sut *DeactivateAvailabilityUsecase
}

func setUpDeactivateAvailability(t *testing.T) *testDeactivateAvailability {
	t.Helper()

	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	consultant, _ := domain.ReconstitueConsultant(
		"con_123",
		"user_123",
		"SOFTWARE_ENGINEER",
		"Jane Doe Tech",
		"Experienced software engineer.",
		10,
		true,
		now,
		now,
	)

	start, _ := domain.NewTimeOfDay(9, 0)
	end, _ := domain.NewTimeOfDay(11, 0)
	// Initial active availability for deactivation testing
	existingAvailability, _ := domain.ReconstitueConsultantAvailability(
		"conav_123456",
		"con_123",
		time.Monday,
		start,
		end,
		true, // active
		now,
		now,
	)

	consultantRepo := &mocks.MockConsultantRepository{
		FindByUserIDFn: func(ctx context.Context, userID string) (*domain.Consultant, error) {
			return consultant, nil
		},
	}

	availabilityRepo := &mocks.MockAvailabilityRepository{
		FindAvailabilityByIDFn: func(ctx context.Context, id string) (*domain.ConsultantAvailability, error) {
			return existingAvailability, nil
		},
		DeactivateAvailabilityFn: func(ctx context.Context, id string) error {
			return nil
		},
	}

	clock := &shared_mocks.MockClock{
		NowFn: func() time.Time {
			return now
		},
	}

	sut := NewDeactivateAvailabilityUsecase(
		availabilityRepo,
		consultantRepo,
		clock,
	)

	return &testDeactivateAvailability{
		consultantRepo:   consultantRepo,
		availabilityRepo: availabilityRepo,
		clock:            clock,
		sut:              sut,
	}
}

func TestDeactivateAvailability_Execute(t *testing.T) {
	userID := "user_123"
	availabilityID := "conav_123456"

	t.Run("should deactivate availability successfully", func(t *testing.T) {
		tc := setUpDeactivateAvailability(t)

		var deactivatedID string
		tc.availabilityRepo.DeactivateAvailabilityFn = func(ctx context.Context, id string) error {
			deactivatedID = id
			return nil
		}

		err := tc.sut.Execute(context.Background(), userID, availabilityID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if deactivatedID != availabilityID {
			t.Errorf("expected deactivated ID %s, got %s", availabilityID, deactivatedID)
		}
	})

	t.Run("should fail when consultant repository returns error", func(t *testing.T) {
		tc := setUpDeactivateAvailability(t)
		expectedErr := errors.New("db query failed")
		tc.consultantRepo.FindByUserIDFn = func(ctx context.Context, userID string) (*domain.Consultant, error) {
			return nil, expectedErr
		}

		err := tc.sut.Execute(context.Background(), userID, availabilityID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("should fail when consultant is not found", func(t *testing.T) {
		tc := setUpDeactivateAvailability(t)
		tc.consultantRepo.FindByUserIDFn = func(ctx context.Context, userID string) (*domain.Consultant, error) {
			return nil, domain.ErrConsultantNotFound
		}

		err := tc.sut.Execute(context.Background(), userID, availabilityID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrConsultantNotFound) {
			t.Errorf("expected ErrConsultantNotFound, got %v", err)
		}
	})

	t.Run("should fail when availability repository FindAvailabilityByID returns error", func(t *testing.T) {
		tc := setUpDeactivateAvailability(t)
		expectedErr := errors.New("availability lookup failed")
		tc.availabilityRepo.FindAvailabilityByIDFn = func(ctx context.Context, id string) (*domain.ConsultantAvailability, error) {
			return nil, expectedErr
		}

		err := tc.sut.Execute(context.Background(), userID, availabilityID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("should fail when availability is not found (nil result)", func(t *testing.T) {
		tc := setUpDeactivateAvailability(t)
		tc.availabilityRepo.FindAvailabilityByIDFn = func(ctx context.Context, id string) (*domain.ConsultantAvailability, error) {
			return nil, nil
		}

		err := tc.sut.Execute(context.Background(), userID, availabilityID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrAvailabilityNotFound) {
			t.Errorf("expected ErrAvailabilityNotFound, got %v", err)
		}
	})

	t.Run("should fail when availability belongs to another consultant", func(t *testing.T) {
		tc := setUpDeactivateAvailability(t)
		otherConsultantAvail, _ := domain.NewConsultantAvailability(
			"conav_123456",
			"con_other_user",
			time.Monday,
			domain.TimeOfDay{},
			domain.TimeOfDay{},
			tc.clock.Now(),
		)
		tc.availabilityRepo.FindAvailabilityByIDFn = func(ctx context.Context, id string) (*domain.ConsultantAvailability, error) {
			return otherConsultantAvail, nil
		}

		err := tc.sut.Execute(context.Background(), userID, availabilityID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrAvailabilityNotFound) {
			t.Errorf("expected ErrAvailabilityNotFound, got %v", err)
		}
	})

	t.Run("should fail when availability is already inactive", func(t *testing.T) {
		tc := setUpDeactivateAvailability(t)
		start, _ := domain.NewTimeOfDay(9, 0)
		end, _ := domain.NewTimeOfDay(11, 0)
		inactiveAvail, _ := domain.ReconstitueConsultantAvailability(
			"conav_123456",
			"con_123",
			time.Monday,
			start,
			end,
			false, // already inactive
			tc.clock.Now(),
			tc.clock.Now(),
		)
		tc.availabilityRepo.FindAvailabilityByIDFn = func(ctx context.Context, id string) (*domain.ConsultantAvailability, error) {
			return inactiveAvail, nil
		}

		err := tc.sut.Execute(context.Background(), userID, availabilityID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrAvailabilityAlreadyDeactivated) {
			t.Errorf("expected ErrAvailabilityAlreadyDeactivated, got %v", err)
		}
	})

	t.Run("should fail when availability repository DeactivateAvailability returns error", func(t *testing.T) {
		tc := setUpDeactivateAvailability(t)
		expectedErr := errors.New("failed to deactivate availability in db")
		tc.availabilityRepo.DeactivateAvailabilityFn = func(ctx context.Context, id string) error {
			return expectedErr
		}

		err := tc.sut.Execute(context.Background(), userID, availabilityID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("should propagate context to repository calls", func(t *testing.T) {
		tc := setUpDeactivateAvailability(t)
		type contextKey string
		key := contextKey("test-ctx-key")
		ctx := context.WithValue(context.Background(), key, "test-ctx-val")

		var receivedFindConsultantCtx context.Context
		origFindByUser := tc.consultantRepo.FindByUserIDFn
		tc.consultantRepo.FindByUserIDFn = func(c context.Context, uID string) (*domain.Consultant, error) {
			receivedFindConsultantCtx = c
			return origFindByUser(c, uID)
		}

		var receivedFindAvailCtx context.Context
		origFindAvail := tc.availabilityRepo.FindAvailabilityByIDFn
		tc.availabilityRepo.FindAvailabilityByIDFn = func(c context.Context, id string) (*domain.ConsultantAvailability, error) {
			receivedFindAvailCtx = c
			return origFindAvail(c, id)
		}

		var receivedDeactivateCtx context.Context
		tc.availabilityRepo.DeactivateAvailabilityFn = func(c context.Context, id string) error {
			receivedDeactivateCtx = c
			return nil
		}

		err := tc.sut.Execute(ctx, userID, availabilityID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if receivedFindConsultantCtx.Value(key) != "test-ctx-val" {
			t.Errorf("expected context in FindByUserID, got %v", receivedFindConsultantCtx.Value(key))
		}
		if receivedFindAvailCtx.Value(key) != "test-ctx-val" {
			t.Errorf("expected context in FindAvailabilityByID, got %v", receivedFindAvailCtx.Value(key))
		}
		if receivedDeactivateCtx.Value(key) != "test-ctx-val" {
			t.Errorf("expected context in DeactivateAvailability, got %v", receivedDeactivateCtx.Value(key))
		}
	})
}
