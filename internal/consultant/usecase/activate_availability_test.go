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

type testActivateAvailability struct {
	consultantRepo   *mocks.MockConsultantRepository
	availabilityRepo *mocks.MockAvailabilityRepository
	clock            *shared_mocks.MockClock

	sut *ActivateAvailabilityUsecase
}

func setUpActivateAvailability(t *testing.T) *testActivateAvailability {
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
	// Initial inactive availability for activation testing
	existingAvailability, _ := domain.ReconstitueConsultantAvailability(
		"conav_123456",
		"con_123",
		time.Monday,
		start,
		end,
		false, // inactive
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
		ActivateAvailabilityFn: func(ctx context.Context, id string) error {
			return nil
		},
	}

	clock := &shared_mocks.MockClock{
		NowFn: func() time.Time {
			return now
		},
	}

	sut := NewActivateAvailabilityUsecase(
		availabilityRepo,
		consultantRepo,
		clock,
	)

	return &testActivateAvailability{
		consultantRepo:   consultantRepo,
		availabilityRepo: availabilityRepo,
		clock:            clock,
		sut:              sut,
	}
}

func TestActivateAvailability_Execute(t *testing.T) {
	userID := "user_123"
	availabilityID := "conav_123456"

	t.Run("should activate availability successfully", func(t *testing.T) {
		tc := setUpActivateAvailability(t)

		var activatedID string
		tc.availabilityRepo.ActivateAvailabilityFn = func(ctx context.Context, id string) error {
			activatedID = id
			return nil
		}

		err := tc.sut.Execute(context.Background(), userID, availabilityID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if activatedID != availabilityID {
			t.Errorf("expected activated ID %s, got %s", availabilityID, activatedID)
		}
	})

	t.Run("should fail when consultant repository returns error", func(t *testing.T) {
		tc := setUpActivateAvailability(t)
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
		tc := setUpActivateAvailability(t)
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
		tc := setUpActivateAvailability(t)
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
		tc := setUpActivateAvailability(t)
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
		tc := setUpActivateAvailability(t)
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

	t.Run("should fail when availability is already active", func(t *testing.T) {
		tc := setUpActivateAvailability(t)
		start, _ := domain.NewTimeOfDay(9, 0)
		end, _ := domain.NewTimeOfDay(11, 0)
		activeAvail, _ := domain.ReconstitueConsultantAvailability(
			"conav_123456",
			"con_123",
			time.Monday,
			start,
			end,
			true, // already active
			tc.clock.Now(),
			tc.clock.Now(),
		)
		tc.availabilityRepo.FindAvailabilityByIDFn = func(ctx context.Context, id string) (*domain.ConsultantAvailability, error) {
			return activeAvail, nil
		}

		err := tc.sut.Execute(context.Background(), userID, availabilityID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrAvailabilityAlreadyActivated) {
			t.Errorf("expected ErrAvailabilityAlreadyActivated, got %v", err)
		}
	})

	t.Run("should fail when availability repository ActivateAvailability returns error", func(t *testing.T) {
		tc := setUpActivateAvailability(t)
		expectedErr := errors.New("failed to activate availability in db")
		tc.availabilityRepo.ActivateAvailabilityFn = func(ctx context.Context, id string) error {
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
		tc := setUpActivateAvailability(t)
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

		var receivedActivateCtx context.Context
		tc.availabilityRepo.ActivateAvailabilityFn = func(c context.Context, id string) error {
			receivedActivateCtx = c
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
		if receivedActivateCtx.Value(key) != "test-ctx-val" {
			t.Errorf("expected context in ActivateAvailability, got %v", receivedActivateCtx.Value(key))
		}
	})
}
