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

type testUpdateAvailability struct {
	consultantRepo   *mocks.MockConsultantRepository
	availabilityRepo *mocks.MockAvailabilityRepository
	clock            *shared_mocks.MockClock

	sut *UpdateAvailabilityUsecase
}

func setUpUpdateAvailability(t *testing.T) *testUpdateAvailability {
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
	existingAvailability, _ := domain.ReconstitueConsultantAvailability(
		"conav_123456",
		"con_123",
		time.Monday,
		start,
		end,
		true,
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
		FindAvailabilitiesByConsultantIDAndDayOfWeekFn: func(ctx context.Context, consultantID string, dayOfWeek time.Weekday) ([]*domain.ConsultantAvailability, error) {
			return []*domain.ConsultantAvailability{existingAvailability}, nil
		},
		UpdateAvailabilityFn: func(ctx context.Context, availability *domain.ConsultantAvailability) error {
			return nil
		},
	}

	clock := &shared_mocks.MockClock{
		NowFn: func() time.Time {
			return now.Add(time.Hour * 24)
		},
	}

	sut := NewUpdateAvailabilityUsecase(
		availabilityRepo,
		consultantRepo,
		clock,
	)

	return &testUpdateAvailability{
		consultantRepo:   consultantRepo,
		availabilityRepo: availabilityRepo,
		clock:            clock,
		sut:              sut,
	}
}

func TestUpdateAvailability_Execute(t *testing.T) {
	updatedTime := time.Date(2026, 1, 2, 12, 0, 0, 0, time.UTC)
	userID := "user_123"
	validReq := &dto.UpdateAvailabilityRequest{
		AvailabilityID: "conav_123456",
		DayOfWeek:      2, // Tuesday
		StartTime:      "14:00",
		EndTime:        "16:00",
	}

	t.Run("should update availability successfully", func(t *testing.T) {
		tc := setUpUpdateAvailability(t)

		var updatedAvailability *domain.ConsultantAvailability
		tc.availabilityRepo.UpdateAvailabilityFn = func(ctx context.Context, availability *domain.ConsultantAvailability) error {
			updatedAvailability = availability
			return nil
		}

		result, err := tc.sut.Execute(context.Background(), userID, validReq)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result == nil {
			t.Fatal("expected result availability, got nil")
		}

		if result.ID() != "conav_123456" {
			t.Errorf("expected ID conav_123456, got %s", result.ID())
		}
		if result.ConsultantID() != "con_123" {
			t.Errorf("expected ConsultantID con_123, got %s", result.ConsultantID())
		}
		if result.DayOfWeek() != time.Tuesday {
			t.Errorf("expected DayOfWeek Tuesday (%v), got %v", time.Tuesday, result.DayOfWeek())
		}
		if result.StartTime().Hour() != 14 || result.StartTime().Minute() != 0 {
			t.Errorf("expected StartTime 14:00, got %02d:%02d", result.StartTime().Hour(), result.StartTime().Minute())
		}
		if result.EndTime().Hour() != 16 || result.EndTime().Minute() != 0 {
			t.Errorf("expected EndTime 16:00, got %02d:%02d", result.EndTime().Hour(), result.EndTime().Minute())
		}
		if !result.IsActive() {
			t.Errorf("expected IsActive to be true")
		}
		if !result.UpdatedAt().Equal(updatedTime) {
			t.Errorf("expected UpdatedAt %v, got %v", updatedTime, result.UpdatedAt())
		}

		if updatedAvailability == nil {
			t.Fatal("expected updated availability sent to repository, got nil")
		}
		if updatedAvailability.ID() != "conav_123456" {
			t.Errorf("expected repo availability ID conav_123456, got %s", updatedAvailability.ID())
		}
	})

	t.Run("should update availability successfully on same day without false self-overlap", func(t *testing.T) {
		tc := setUpUpdateAvailability(t)

		req := &dto.UpdateAvailabilityRequest{
			AvailabilityID: "conav_123456",
			DayOfWeek:      1, // Monday (same day)
			StartTime:      "10:00",
			EndTime:        "12:00",
		}

		result, err := tc.sut.Execute(context.Background(), userID, req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result.DayOfWeek() != time.Monday {
			t.Errorf("expected DayOfWeek Monday, got %v", result.DayOfWeek())
		}
		if result.StartTime().String() != "10:00" {
			t.Errorf("expected start time 10:00, got %s", result.StartTime().String())
		}
		if result.EndTime().String() != "12:00" {
			t.Errorf("expected end time 12:00, got %s", result.EndTime().String())
		}
	})

	t.Run("should update availability successfully when other slot on same day is adjacent before", func(t *testing.T) {
		tc := setUpUpdateAvailability(t)

		otherStart, _ := domain.NewTimeOfDay(8, 0)
		otherEnd, _ := domain.NewTimeOfDay(10, 0)
		otherSlot, _ := domain.NewConsultantAvailability(
			"conav_other",
			"con_123",
			time.Monday,
			otherStart,
			otherEnd,
			tc.clock.Now(),
		)

		currentSlot, _ := domain.NewConsultantAvailability(
			"conav_123456",
			"con_123",
			time.Monday,
			otherStart,
			otherEnd,
			tc.clock.Now(),
		)

		tc.availabilityRepo.FindAvailabilityByIDFn = func(ctx context.Context, id string) (*domain.ConsultantAvailability, error) {
			return currentSlot, nil
		}
		tc.availabilityRepo.FindAvailabilitiesByConsultantIDAndDayOfWeekFn = func(ctx context.Context, consultantID string, dayOfWeek time.Weekday) ([]*domain.ConsultantAvailability, error) {
			return []*domain.ConsultantAvailability{currentSlot, otherSlot}, nil
		}

		// Updating currentSlot to 10:00 - 12:00 (adjacent after otherSlot 08:00 - 10:00)
		req := &dto.UpdateAvailabilityRequest{
			AvailabilityID: "conav_123456",
			DayOfWeek:      1,
			StartTime:      "10:00",
			EndTime:        "12:00",
		}

		result, err := tc.sut.Execute(context.Background(), userID, req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.StartTime().String() != "10:00" || result.EndTime().String() != "12:00" {
			t.Errorf("expected 10:00-12:00, got %s-%s", result.StartTime().String(), result.EndTime().String())
		}
	})

	t.Run("should update availability successfully when other slot on same day is adjacent after", func(t *testing.T) {
		tc := setUpUpdateAvailability(t)

		otherStart, _ := domain.NewTimeOfDay(12, 0)
		otherEnd, _ := domain.NewTimeOfDay(14, 0)
		otherSlot, _ := domain.NewConsultantAvailability(
			"conav_other",
			"con_123",
			time.Monday,
			otherStart,
			otherEnd,
			tc.clock.Now(),
		)

		currentSlot, _ := domain.NewConsultantAvailability(
			"conav_123456",
			"con_123",
			time.Monday,
			otherStart,
			otherEnd,
			tc.clock.Now(),
		)

		tc.availabilityRepo.FindAvailabilityByIDFn = func(ctx context.Context, id string) (*domain.ConsultantAvailability, error) {
			return currentSlot, nil
		}
		tc.availabilityRepo.FindAvailabilitiesByConsultantIDAndDayOfWeekFn = func(ctx context.Context, consultantID string, dayOfWeek time.Weekday) ([]*domain.ConsultantAvailability, error) {
			return []*domain.ConsultantAvailability{currentSlot, otherSlot}, nil
		}

		// Updating currentSlot to 10:00 - 12:00 (adjacent before otherSlot 12:00 - 14:00)
		req := &dto.UpdateAvailabilityRequest{
			AvailabilityID: "conav_123456",
			DayOfWeek:      1,
			StartTime:      "10:00",
			EndTime:        "12:00",
		}

		result, err := tc.sut.Execute(context.Background(), userID, req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.StartTime().String() != "10:00" || result.EndTime().String() != "12:00" {
			t.Errorf("expected 10:00-12:00, got %s-%s", result.StartTime().String(), result.EndTime().String())
		}
	})

	t.Run("should fail when consultant repository returns error", func(t *testing.T) {
		tc := setUpUpdateAvailability(t)
		expectedErr := errors.New("database connection failed")
		tc.consultantRepo.FindByUserIDFn = func(ctx context.Context, userID string) (*domain.Consultant, error) {
			return nil, expectedErr
		}

		_, err := tc.sut.Execute(context.Background(), userID, validReq)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("should fail when consultant is not found", func(t *testing.T) {
		tc := setUpUpdateAvailability(t)
		tc.consultantRepo.FindByUserIDFn = func(ctx context.Context, userID string) (*domain.Consultant, error) {
			return nil, domain.ErrConsultantNotFound
		}

		_, err := tc.sut.Execute(context.Background(), userID, validReq)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrConsultantNotFound) {
			t.Errorf("expected ErrConsultantNotFound, got %v", err)
		}
	})

	t.Run("should fail when availability repository FindAvailabilityByID returns error", func(t *testing.T) {
		tc := setUpUpdateAvailability(t)
		expectedErr := errors.New("failed to query availability")
		tc.availabilityRepo.FindAvailabilityByIDFn = func(ctx context.Context, id string) (*domain.ConsultantAvailability, error) {
			return nil, expectedErr
		}

		_, err := tc.sut.Execute(context.Background(), userID, validReq)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("should fail when availability is not found (nil result)", func(t *testing.T) {
		tc := setUpUpdateAvailability(t)
		tc.availabilityRepo.FindAvailabilityByIDFn = func(ctx context.Context, id string) (*domain.ConsultantAvailability, error) {
			return nil, nil
		}

		_, err := tc.sut.Execute(context.Background(), userID, validReq)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrAvailabilityNotFound) {
			t.Errorf("expected ErrAvailabilityNotFound, got %v", err)
		}
	})

	t.Run("should fail when availability belongs to another consultant", func(t *testing.T) {
		tc := setUpUpdateAvailability(t)
		otherConsultantAvail, _ := domain.NewConsultantAvailability(
			"conav_123456",
			"con_other_999", // Different consultant
			time.Monday,
			domain.TimeOfDay{},
			domain.TimeOfDay{},
			tc.clock.Now(),
		)
		tc.availabilityRepo.FindAvailabilityByIDFn = func(ctx context.Context, id string) (*domain.ConsultantAvailability, error) {
			return otherConsultantAvail, nil
		}

		_, err := tc.sut.Execute(context.Background(), userID, validReq)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrAvailabilityNotFound) {
			t.Errorf("expected ErrAvailabilityNotFound, got %v", err)
		}
	})

	t.Run("should fail when start time format is invalid", func(t *testing.T) {
		tc := setUpUpdateAvailability(t)
		req := &dto.UpdateAvailabilityRequest{
			AvailabilityID: "conav_123456",
			DayOfWeek:      2,
			StartTime:      "invalid",
			EndTime:        "16:00",
		}

		_, err := tc.sut.Execute(context.Background(), userID, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrInvalidTimeFormat) {
			t.Errorf("expected ErrInvalidTimeFormat, got %v", err)
		}
	})

	t.Run("should fail when start time hour is invalid", func(t *testing.T) {
		tc := setUpUpdateAvailability(t)
		req := &dto.UpdateAvailabilityRequest{
			AvailabilityID: "conav_123456",
			DayOfWeek:      2,
			StartTime:      "25:00",
			EndTime:        "16:00",
		}

		_, err := tc.sut.Execute(context.Background(), userID, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrInvalidHour) {
			t.Errorf("expected ErrInvalidHour, got %v", err)
		}
	})

	t.Run("should fail when start time minute is invalid", func(t *testing.T) {
		tc := setUpUpdateAvailability(t)
		req := &dto.UpdateAvailabilityRequest{
			AvailabilityID: "conav_123456",
			DayOfWeek:      2,
			StartTime:      "14:60",
			EndTime:        "16:00",
		}

		_, err := tc.sut.Execute(context.Background(), userID, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrInvalidMinute) {
			t.Errorf("expected ErrInvalidMinute, got %v", err)
		}
	})

	t.Run("should fail when end time format is invalid", func(t *testing.T) {
		tc := setUpUpdateAvailability(t)
		req := &dto.UpdateAvailabilityRequest{
			AvailabilityID: "conav_123456",
			DayOfWeek:      2,
			StartTime:      "14:00",
			EndTime:        "invalid",
		}

		_, err := tc.sut.Execute(context.Background(), userID, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrInvalidTimeFormat) {
			t.Errorf("expected ErrInvalidTimeFormat, got %v", err)
		}
	})

	t.Run("should fail when end time hour is invalid", func(t *testing.T) {
		tc := setUpUpdateAvailability(t)
		req := &dto.UpdateAvailabilityRequest{
			AvailabilityID: "conav_123456",
			DayOfWeek:      2,
			StartTime:      "14:00",
			EndTime:        "24:00",
		}

		_, err := tc.sut.Execute(context.Background(), userID, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrInvalidHour) {
			t.Errorf("expected ErrInvalidHour, got %v", err)
		}
	})

	t.Run("should fail when end time minute is invalid", func(t *testing.T) {
		tc := setUpUpdateAvailability(t)
		req := &dto.UpdateAvailabilityRequest{
			AvailabilityID: "conav_123456",
			DayOfWeek:      2,
			StartTime:      "14:00",
			EndTime:        "15:70",
		}

		_, err := tc.sut.Execute(context.Background(), userID, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrInvalidMinute) {
			t.Errorf("expected ErrInvalidMinute, got %v", err)
		}
	})

	t.Run("should fail when start time is equal to end time", func(t *testing.T) {
		tc := setUpUpdateAvailability(t)
		req := &dto.UpdateAvailabilityRequest{
			AvailabilityID: "conav_123456",
			DayOfWeek:      2,
			StartTime:      "14:00",
			EndTime:        "14:00",
		}

		_, err := tc.sut.Execute(context.Background(), userID, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrInvalidTimeRange) {
			t.Errorf("expected ErrInvalidTimeRange, got %v", err)
		}
	})

	t.Run("should fail when start time is after end time", func(t *testing.T) {
		tc := setUpUpdateAvailability(t)
		req := &dto.UpdateAvailabilityRequest{
			AvailabilityID: "conav_123456",
			DayOfWeek:      2,
			StartTime:      "16:00",
			EndTime:        "14:00",
		}

		_, err := tc.sut.Execute(context.Background(), userID, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrInvalidTimeRange) {
			t.Errorf("expected ErrInvalidTimeRange, got %v", err)
		}
	})

	t.Run("should fail when availability repository FindAvailabilitiesByConsultantIDAndDayOfWeek returns error", func(t *testing.T) {
		tc := setUpUpdateAvailability(t)
		expectedErr := errors.New("failed to query day availabilities")
		tc.availabilityRepo.FindAvailabilitiesByConsultantIDAndDayOfWeekFn = func(ctx context.Context, consultantID string, dayOfWeek time.Weekday) ([]*domain.ConsultantAvailability, error) {
			return nil, expectedErr
		}

		_, err := tc.sut.Execute(context.Background(), userID, validReq)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("should fail when updated availability overlaps with another slot (exact match)", func(t *testing.T) {
		tc := setUpUpdateAvailability(t)

		otherStart, _ := domain.NewTimeOfDay(14, 0)
		otherEnd, _ := domain.NewTimeOfDay(16, 0)
		otherSlot, _ := domain.NewConsultantAvailability(
			"conav_other_slot",
			"con_123",
			time.Tuesday,
			otherStart,
			otherEnd,
			tc.clock.Now(),
		)

		tc.availabilityRepo.FindAvailabilitiesByConsultantIDAndDayOfWeekFn = func(ctx context.Context, consultantID string, dayOfWeek time.Weekday) ([]*domain.ConsultantAvailability, error) {
			return []*domain.ConsultantAvailability{otherSlot}, nil
		}

		_, err := tc.sut.Execute(context.Background(), userID, validReq)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrAvailabilityOverlap) {
			t.Errorf("expected ErrAvailabilityOverlap, got %v", err)
		}
	})

	t.Run("should fail when updated availability starts inside another slot", func(t *testing.T) {
		tc := setUpUpdateAvailability(t)

		otherStart, _ := domain.NewTimeOfDay(13, 0)
		otherEnd, _ := domain.NewTimeOfDay(15, 0)
		otherSlot, _ := domain.NewConsultantAvailability(
			"conav_other_slot",
			"con_123",
			time.Tuesday,
			otherStart,
			otherEnd,
			tc.clock.Now(),
		)

		tc.availabilityRepo.FindAvailabilitiesByConsultantIDAndDayOfWeekFn = func(ctx context.Context, consultantID string, dayOfWeek time.Weekday) ([]*domain.ConsultantAvailability, error) {
			return []*domain.ConsultantAvailability{otherSlot}, nil
		}

		_, err := tc.sut.Execute(context.Background(), userID, validReq)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrAvailabilityOverlap) {
			t.Errorf("expected ErrAvailabilityOverlap, got %v", err)
		}
	})

	t.Run("should fail when updated availability ends inside another slot", func(t *testing.T) {
		tc := setUpUpdateAvailability(t)

		otherStart, _ := domain.NewTimeOfDay(15, 0)
		otherEnd, _ := domain.NewTimeOfDay(17, 0)
		otherSlot, _ := domain.NewConsultantAvailability(
			"conav_other_slot",
			"con_123",
			time.Tuesday,
			otherStart,
			otherEnd,
			tc.clock.Now(),
		)

		tc.availabilityRepo.FindAvailabilitiesByConsultantIDAndDayOfWeekFn = func(ctx context.Context, consultantID string, dayOfWeek time.Weekday) ([]*domain.ConsultantAvailability, error) {
			return []*domain.ConsultantAvailability{otherSlot}, nil
		}

		_, err := tc.sut.Execute(context.Background(), userID, validReq)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrAvailabilityOverlap) {
			t.Errorf("expected ErrAvailabilityOverlap, got %v", err)
		}
	})

	t.Run("should fail when updated availability completely engulfs another slot", func(t *testing.T) {
		tc := setUpUpdateAvailability(t)

		otherStart, _ := domain.NewTimeOfDay(14, 30)
		otherEnd, _ := domain.NewTimeOfDay(15, 30)
		otherSlot, _ := domain.NewConsultantAvailability(
			"conav_other_slot",
			"con_123",
			time.Tuesday,
			otherStart,
			otherEnd,
			tc.clock.Now(),
		)

		tc.availabilityRepo.FindAvailabilitiesByConsultantIDAndDayOfWeekFn = func(ctx context.Context, consultantID string, dayOfWeek time.Weekday) ([]*domain.ConsultantAvailability, error) {
			return []*domain.ConsultantAvailability{otherSlot}, nil
		}

		_, err := tc.sut.Execute(context.Background(), userID, validReq)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrAvailabilityOverlap) {
			t.Errorf("expected ErrAvailabilityOverlap, got %v", err)
		}
	})

	t.Run("should fail when availability repository UpdateAvailability returns error", func(t *testing.T) {
		tc := setUpUpdateAvailability(t)
		expectedErr := errors.New("failed to persist availability update")
		tc.availabilityRepo.UpdateAvailabilityFn = func(ctx context.Context, availability *domain.ConsultantAvailability) error {
			return expectedErr
		}

		_, err := tc.sut.Execute(context.Background(), userID, validReq)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("should propagate context to repository calls", func(t *testing.T) {
		tc := setUpUpdateAvailability(t)
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

		var receivedFindAvailsByDayCtx context.Context
		origFindAvailsByDay := tc.availabilityRepo.FindAvailabilitiesByConsultantIDAndDayOfWeekFn
		tc.availabilityRepo.FindAvailabilitiesByConsultantIDAndDayOfWeekFn = func(c context.Context, cID string, dayOfWeek time.Weekday) ([]*domain.ConsultantAvailability, error) {
			receivedFindAvailsByDayCtx = c
			return origFindAvailsByDay(c, cID, dayOfWeek)
		}

		var receivedUpdateCtx context.Context
		tc.availabilityRepo.UpdateAvailabilityFn = func(c context.Context, availability *domain.ConsultantAvailability) error {
			receivedUpdateCtx = c
			return nil
		}

		_, err := tc.sut.Execute(ctx, userID, validReq)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if receivedFindConsultantCtx.Value(key) != "test-ctx-val" {
			t.Errorf("expected context in FindByUserID, got %v", receivedFindConsultantCtx.Value(key))
		}
		if receivedFindAvailCtx.Value(key) != "test-ctx-val" {
			t.Errorf("expected context in FindAvailabilityByID, got %v", receivedFindAvailCtx.Value(key))
		}
		if receivedFindAvailsByDayCtx.Value(key) != "test-ctx-val" {
			t.Errorf("expected context in FindAvailabilitiesByConsultantIDAndDayOfWeek, got %v", receivedFindAvailsByDayCtx.Value(key))
		}
		if receivedUpdateCtx.Value(key) != "test-ctx-val" {
			t.Errorf("expected context in UpdateAvailability, got %v", receivedUpdateCtx.Value(key))
		}
	})
}
