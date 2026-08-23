package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultant/mocks"
)

type testGetAvailability struct {
	availabilityRepo *mocks.MockAvailabilityRepository
	sut              *GetAvailabilityUsecase
}

func setUpGetAvailability(t *testing.T) *testGetAvailability {
	t.Helper()

	availabilityRepo := &mocks.MockAvailabilityRepository{}
	sut := NewGetAvailabilityUsecase(availabilityRepo)

	return &testGetAvailability{
		availabilityRepo: availabilityRepo,
		sut:              sut,
	}
}

func TestGetAvailability_Execute(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	consultantID := "con_123"

	start1, err := domain.NewTimeOfDay(9, 0)
	if err != nil {
		t.Fatalf("failed to create start time: %v", err)
	}
	end1, err := domain.NewTimeOfDay(11, 0)
	if err != nil {
		t.Fatalf("failed to create end time: %v", err)
	}
	avail1, err := domain.ReconstitueConsultantAvailability(
		"conav_001",
		consultantID,
		time.Monday,
		start1,
		end1,
		true,
		now,
		now,
	)
	if err != nil {
		t.Fatalf("failed to create availability 1: %v", err)
	}

	start2, err := domain.NewTimeOfDay(14, 0)
	if err != nil {
		t.Fatalf("failed to create start time: %v", err)
	}
	end2, err := domain.NewTimeOfDay(16, 0)
	if err != nil {
		t.Fatalf("failed to create end time: %v", err)
	}
	avail2, err := domain.ReconstitueConsultantAvailability(
		"conav_002",
		consultantID,
		time.Tuesday,
		start2,
		end2,
		true,
		now,
		now,
	)
	if err != nil {
		t.Fatalf("failed to create availability 2: %v", err)
	}

	t.Run("should get availabilities for consultant successfully when availabilities exist", func(t *testing.T) {
		tc := setUpGetAvailability(t)
		tc.availabilityRepo.FindAvailabilitiesByConsultantIDFn = func(ctx context.Context, cID string) ([]*domain.ConsultantAvailability, error) {
			if cID != consultantID {
				t.Errorf("expected consultantID %s, got %s", consultantID, cID)
			}
			return []*domain.ConsultantAvailability{avail1, avail2}, nil
		}

		result, err := tc.sut.Execute(context.Background(), consultantID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(result) != 2 {
			t.Fatalf("expected 2 availabilities, got %d", len(result))
		}

		// Verify first slot
		if result[0].ID() != "conav_001" {
			t.Errorf("expected ID conav_001, got %s", result[0].ID())
		}
		if result[0].ConsultantID() != consultantID {
			t.Errorf("expected consultantID %s, got %s", consultantID, result[0].ConsultantID())
		}
		if result[0].DayOfWeek() != time.Monday {
			t.Errorf("expected DayOfWeek Monday (1), got %v", result[0].DayOfWeek())
		}
		if result[0].StartTime().String() != "09:00" {
			t.Errorf("expected start time 09:00, got %s", result[0].StartTime().String())
		}
		if result[0].EndTime().String() != "11:00" {
			t.Errorf("expected end time 11:00, got %s", result[0].EndTime().String())
		}
		if !result[0].IsActive() {
			t.Errorf("expected isActive true, got %v", result[0].IsActive())
		}

		// Verify second slot
		if result[1].ID() != "conav_002" {
			t.Errorf("expected ID conav_002, got %s", result[1].ID())
		}
		if result[1].ConsultantID() != consultantID {
			t.Errorf("expected consultantID %s, got %s", consultantID, result[1].ConsultantID())
		}
		if result[1].DayOfWeek() != time.Tuesday {
			t.Errorf("expected DayOfWeek Tuesday (2), got %v", result[1].DayOfWeek())
		}
		if result[1].StartTime().String() != "14:00" {
			t.Errorf("expected start time 14:00, got %s", result[1].StartTime().String())
		}
		if result[1].EndTime().String() != "16:00" {
			t.Errorf("expected end time 16:00, got %s", result[1].EndTime().String())
		}
		if !result[1].IsActive() {
			t.Errorf("expected isActive true, got %v", result[1].IsActive())
		}
	})

	t.Run("should return empty list when no availabilities exist for consultant", func(t *testing.T) {
		tc := setUpGetAvailability(t)
		tc.availabilityRepo.FindAvailabilitiesByConsultantIDFn = func(ctx context.Context, cID string) ([]*domain.ConsultantAvailability, error) {
			return []*domain.ConsultantAvailability{}, nil
		}

		result, err := tc.sut.Execute(context.Background(), consultantID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(result) != 0 {
			t.Fatalf("expected 0 availabilities, got %d", len(result))
		}
	})

	t.Run("should fail when availability repository FindAvailabilitiesByConsultantID returns error", func(t *testing.T) {
		tc := setUpGetAvailability(t)
		expectedErr := errors.New("db query failed")
		tc.availabilityRepo.FindAvailabilitiesByConsultantIDFn = func(ctx context.Context, cID string) ([]*domain.ConsultantAvailability, error) {
			return nil, expectedErr
		}

		result, err := tc.sut.Execute(context.Background(), consultantID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
		if result != nil {
			t.Errorf("expected nil result, got %v", result)
		}
	})

	t.Run("should propagate context to repository", func(t *testing.T) {
		tc := setUpGetAvailability(t)
		type contextKey string
		key := contextKey("test-key")
		ctx := context.WithValue(context.Background(), key, "test-value")

		var receivedCtx context.Context
		tc.availabilityRepo.FindAvailabilitiesByConsultantIDFn = func(c context.Context, cID string) ([]*domain.ConsultantAvailability, error) {
			receivedCtx = c
			return []*domain.ConsultantAvailability{}, nil
		}

		_, err := tc.sut.Execute(ctx, consultantID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if receivedCtx.Value(key) != "test-value" {
			t.Errorf("expected context value test-value, got %v", receivedCtx.Value(key))
		}
	})
}
