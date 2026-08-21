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

type testCreateAvailability struct {
	consultantRepo   *mocks.MockConsultantRepository
	availabilityRepo *mocks.MockAvailabilityRepository
	idGenerator      *shared_mocks.MockIDGenerator
	clock            *shared_mocks.MockClock

	sut *CreateAvailabilityUsecase
}

func setUpCreateAvailability(t *testing.T) *testCreateAvailability {
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

	consultantRepo := &mocks.MockConsultantRepository{
		FindByUserIDFn: func(ctx context.Context, userID string) (*domain.Consultant, error) {
			return consultant, nil
		},
	}

	availabilityRepo := &mocks.MockAvailabilityRepository{
		FindAvailabilitiesByConsultantIDAndDayOfWeekFn: func(ctx context.Context, consultantID string, dayOfWeek time.Weekday) ([]*domain.ConsultantAvailability, error) {
			return []*domain.ConsultantAvailability{}, nil
		},
		SaveAvailabilityFn: func(ctx context.Context, availability *domain.ConsultantAvailability) error {
			return nil
		},
	}

	idGenerator := &shared_mocks.MockIDGenerator{
		GenerateFn: func(prefix string) (string, error) {
			return "conav_123456", nil
		},
	}

	clock := &shared_mocks.MockClock{
		NowFn: func() time.Time {
			return now
		},
	}

	sut := NewCreateAvailabilityUsecase(
		availabilityRepo,
		idGenerator,
		consultantRepo,
		clock,
	)

	return &testCreateAvailability{
		consultantRepo:   consultantRepo,
		availabilityRepo: availabilityRepo,
		idGenerator:      idGenerator,
		clock:            clock,
		sut:              sut,
	}
}

func TestCreateAvailability_Execute(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	userID := "user_123"
	validReq := &dto.CreateAvailabilityRequest{
		DayOfWeek: 1, // Monday
		StartTime: "09:00",
		EndTime:   "11:00",
	}

	t.Run("should create availability successfully", func(t *testing.T) {
		tc := setUpCreateAvailability(t)

		var generatedPrefix string
		tc.idGenerator.GenerateFn = func(prefix string) (string, error) {
			generatedPrefix = prefix
			return "conav_123456", nil
		}

		var savedAvailability *domain.ConsultantAvailability
		tc.availabilityRepo.SaveAvailabilityFn = func(ctx context.Context, availability *domain.ConsultantAvailability) error {
			savedAvailability = availability
			return nil
		}

		err := tc.sut.Execute(context.Background(), userID, validReq)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if generatedPrefix != domain.ConsultantAvailabilityIDPrefix {
			t.Errorf("expected ID prefix %s, got %s", domain.ConsultantAvailabilityIDPrefix, generatedPrefix)
		}

		if savedAvailability == nil {
			t.Fatal("expected saved availability, got nil")
		}

		if savedAvailability.ID() != "conav_123456" {
			t.Errorf("expected ID conav_123456, got %s", savedAvailability.ID())
		}
		if savedAvailability.ConsultantID() != "con_123" {
			t.Errorf("expected ConsultantID con_123, got %s", savedAvailability.ConsultantID())
		}
		if savedAvailability.DayOfWeek() != time.Monday {
			t.Errorf("expected DayOfWeek Monday (%v), got %v", time.Monday, savedAvailability.DayOfWeek())
		}
		if savedAvailability.StartTime().Hour() != 9 || savedAvailability.StartTime().Minute() != 0 {
			t.Errorf("expected StartTime 09:00, got %02d:%02d", savedAvailability.StartTime().Hour(), savedAvailability.StartTime().Minute())
		}
		if savedAvailability.EndTime().Hour() != 11 || savedAvailability.EndTime().Minute() != 0 {
			t.Errorf("expected EndTime 11:00, got %02d:%02d", savedAvailability.EndTime().Hour(), savedAvailability.EndTime().Minute())
		}
		if !savedAvailability.IsActive() {
			t.Errorf("expected IsActive to be true")
		}
		if !savedAvailability.CreatedAt().Equal(now) {
			t.Errorf("expected CreatedAt %v, got %v", now, savedAvailability.CreatedAt())
		}
		if !savedAvailability.UpdatedAt().Equal(now) {
			t.Errorf("expected UpdatedAt %v, got %v", now, savedAvailability.UpdatedAt())
		}
	})

	t.Run("should create availability successfully when existing availability is adjacent before", func(t *testing.T) {
		tc := setUpCreateAvailability(t)

		existingStart, _ := domain.NewTimeOfDay(11, 0)
		existingEnd, _ := domain.NewTimeOfDay(13, 0)
		existingAvail, _ := domain.NewConsultantAvailability(
			"conav_existing",
			"con_123",
			time.Monday,
			existingStart,
			existingEnd,
			tc.clock.Now(),
		)

		tc.availabilityRepo.FindAvailabilitiesByConsultantIDAndDayOfWeekFn = func(ctx context.Context, consultantID string, dayOfWeek time.Weekday) ([]*domain.ConsultantAvailability, error) {
			return []*domain.ConsultantAvailability{existingAvail}, nil
		}

		// New: 09:00 - 11:00 (adjacent before existing 11:00 - 13:00)
		err := tc.sut.Execute(context.Background(), userID, validReq)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("should create availability successfully when existing availability is adjacent after", func(t *testing.T) {
		tc := setUpCreateAvailability(t)

		existingStart, _ := domain.NewTimeOfDay(7, 0)
		existingEnd, _ := domain.NewTimeOfDay(9, 0)
		existingAvail, _ := domain.NewConsultantAvailability(
			"conav_existing",
			"con_123",
			time.Monday,
			existingStart,
			existingEnd,
			tc.clock.Now(),
		)

		tc.availabilityRepo.FindAvailabilitiesByConsultantIDAndDayOfWeekFn = func(ctx context.Context, consultantID string, dayOfWeek time.Weekday) ([]*domain.ConsultantAvailability, error) {
			return []*domain.ConsultantAvailability{existingAvail}, nil
		}

		// New: 09:00 - 11:00 (adjacent after existing 07:00 - 09:00)
		err := tc.sut.Execute(context.Background(), userID, validReq)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("should fail when consultant repository returns error", func(t *testing.T) {
		tc := setUpCreateAvailability(t)
		expectedErr := errors.New("database connection failed")
		tc.consultantRepo.FindByUserIDFn = func(ctx context.Context, userID string) (*domain.Consultant, error) {
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

	t.Run("should fail when consultant is not found", func(t *testing.T) {
		tc := setUpCreateAvailability(t)
		tc.consultantRepo.FindByUserIDFn = func(ctx context.Context, userID string) (*domain.Consultant, error) {
			return nil, domain.ErrConsultantNotFound
		}

		err := tc.sut.Execute(context.Background(), userID, validReq)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrConsultantNotFound) {
			t.Errorf("expected ErrConsultantNotFound, got %v", err)
		}
	})

	t.Run("should fail when consultant is not accepting clients", func(t *testing.T) {
		tc := setUpCreateAvailability(t)
		inactiveConsultant, _ := domain.ReconstitueConsultant(
			"con_123",
			"user_123",
			"SOFTWARE_ENGINEER",
			"Jane Doe Tech",
			"Experienced software engineer.",
			10,
			false, // not accepting clients
			now,
			now,
		)
		tc.consultantRepo.FindByUserIDFn = func(ctx context.Context, userID string) (*domain.Consultant, error) {
			return inactiveConsultant, nil
		}

		err := tc.sut.Execute(context.Background(), userID, validReq)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrConsultantNotAcceptingClients) {
			t.Errorf("expected ErrConsultantNotAcceptingClients, got %v", err)
		}
	})

	t.Run("should fail when start time format is invalid", func(t *testing.T) {
		tc := setUpCreateAvailability(t)
		req := &dto.CreateAvailabilityRequest{
			DayOfWeek: 1,
			StartTime: "invalid",
			EndTime:   "11:00",
		}

		err := tc.sut.Execute(context.Background(), userID, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrInvalidTimeFormat) {
			t.Errorf("expected ErrInvalidTimeFormat, got %v", err)
		}
	})

	t.Run("should fail when start time hour is invalid", func(t *testing.T) {
		tc := setUpCreateAvailability(t)
		req := &dto.CreateAvailabilityRequest{
			DayOfWeek: 1,
			StartTime: "25:00",
			EndTime:   "11:00",
		}

		err := tc.sut.Execute(context.Background(), userID, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrInvalidHour) {
			t.Errorf("expected ErrInvalidHour, got %v", err)
		}
	})

	t.Run("should fail when start time minute is invalid", func(t *testing.T) {
		tc := setUpCreateAvailability(t)
		req := &dto.CreateAvailabilityRequest{
			DayOfWeek: 1,
			StartTime: "09:60",
			EndTime:   "11:00",
		}

		err := tc.sut.Execute(context.Background(), userID, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrInvalidMinute) {
			t.Errorf("expected ErrInvalidMinute, got %v", err)
		}
	})

	t.Run("should fail when end time format is invalid", func(t *testing.T) {
		tc := setUpCreateAvailability(t)
		req := &dto.CreateAvailabilityRequest{
			DayOfWeek: 1,
			StartTime: "09:00",
			EndTime:   "invalid",
		}

		err := tc.sut.Execute(context.Background(), userID, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrInvalidTimeFormat) {
			t.Errorf("expected ErrInvalidTimeFormat, got %v", err)
		}
	})

	t.Run("should fail when end time hour is invalid", func(t *testing.T) {
		tc := setUpCreateAvailability(t)
		req := &dto.CreateAvailabilityRequest{
			DayOfWeek: 1,
			StartTime: "09:00",
			EndTime:   "24:00",
		}

		err := tc.sut.Execute(context.Background(), userID, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrInvalidHour) {
			t.Errorf("expected ErrInvalidHour, got %v", err)
		}
	})

	t.Run("should fail when end time minute is invalid", func(t *testing.T) {
		tc := setUpCreateAvailability(t)
		req := &dto.CreateAvailabilityRequest{
			DayOfWeek: 1,
			StartTime: "09:00",
			EndTime:   "10:75",
		}

		err := tc.sut.Execute(context.Background(), userID, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrInvalidMinute) {
			t.Errorf("expected ErrInvalidMinute, got %v", err)
		}
	})

	t.Run("should fail when start time is equal to end time", func(t *testing.T) {
		tc := setUpCreateAvailability(t)
		req := &dto.CreateAvailabilityRequest{
			DayOfWeek: 1,
			StartTime: "10:00",
			EndTime:   "10:00",
		}

		err := tc.sut.Execute(context.Background(), userID, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrInvalidTimeRange) {
			t.Errorf("expected ErrInvalidTimeRange, got %v", err)
		}
	})

	t.Run("should fail when start time is after end time", func(t *testing.T) {
		tc := setUpCreateAvailability(t)
		req := &dto.CreateAvailabilityRequest{
			DayOfWeek: 1,
			StartTime: "11:00",
			EndTime:   "09:00",
		}

		err := tc.sut.Execute(context.Background(), userID, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrInvalidTimeRange) {
			t.Errorf("expected ErrInvalidTimeRange, got %v", err)
		}
	})

	t.Run("should fail when availability repository FindAvailabilities returns error", func(t *testing.T) {
		tc := setUpCreateAvailability(t)
		expectedErr := errors.New("failed to query availabilities")
		tc.availabilityRepo.FindAvailabilitiesByConsultantIDAndDayOfWeekFn = func(ctx context.Context, consultantID string, dayOfWeek time.Weekday) ([]*domain.ConsultantAvailability, error) {
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

	t.Run("should fail when new availability overlaps with existing availability (exact match)", func(t *testing.T) {
		tc := setUpCreateAvailability(t)

		existingStart, _ := domain.NewTimeOfDay(9, 0)
		existingEnd, _ := domain.NewTimeOfDay(11, 0)
		existingAvail, _ := domain.NewConsultantAvailability(
			"conav_existing",
			"con_123",
			time.Monday,
			existingStart,
			existingEnd,
			tc.clock.Now(),
		)

		tc.availabilityRepo.FindAvailabilitiesByConsultantIDAndDayOfWeekFn = func(ctx context.Context, consultantID string, dayOfWeek time.Weekday) ([]*domain.ConsultantAvailability, error) {
			return []*domain.ConsultantAvailability{existingAvail}, nil
		}

		err := tc.sut.Execute(context.Background(), userID, validReq)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrAvailabilityOverlap) {
			t.Errorf("expected ErrAvailabilityOverlap, got %v", err)
		}
	})

	t.Run("should fail when new availability starts inside existing availability", func(t *testing.T) {
		tc := setUpCreateAvailability(t)

		existingStart, _ := domain.NewTimeOfDay(9, 0)
		existingEnd, _ := domain.NewTimeOfDay(11, 0)
		existingAvail, _ := domain.NewConsultantAvailability(
			"conav_existing",
			"con_123",
			time.Monday,
			existingStart,
			existingEnd,
			tc.clock.Now(),
		)

		tc.availabilityRepo.FindAvailabilitiesByConsultantIDAndDayOfWeekFn = func(ctx context.Context, consultantID string, dayOfWeek time.Weekday) ([]*domain.ConsultantAvailability, error) {
			return []*domain.ConsultantAvailability{existingAvail}, nil
		}

		req := &dto.CreateAvailabilityRequest{
			DayOfWeek: 1,
			StartTime: "10:00",
			EndTime:   "12:00",
		}

		err := tc.sut.Execute(context.Background(), userID, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrAvailabilityOverlap) {
			t.Errorf("expected ErrAvailabilityOverlap, got %v", err)
		}
	})

	t.Run("should fail when new availability ends inside existing availability", func(t *testing.T) {
		tc := setUpCreateAvailability(t)

		existingStart, _ := domain.NewTimeOfDay(9, 0)
		existingEnd, _ := domain.NewTimeOfDay(11, 0)
		existingAvail, _ := domain.NewConsultantAvailability(
			"conav_existing",
			"con_123",
			time.Monday,
			existingStart,
			existingEnd,
			tc.clock.Now(),
		)

		tc.availabilityRepo.FindAvailabilitiesByConsultantIDAndDayOfWeekFn = func(ctx context.Context, consultantID string, dayOfWeek time.Weekday) ([]*domain.ConsultantAvailability, error) {
			return []*domain.ConsultantAvailability{existingAvail}, nil
		}

		req := &dto.CreateAvailabilityRequest{
			DayOfWeek: 1,
			StartTime: "08:00",
			EndTime:   "10:00",
		}

		err := tc.sut.Execute(context.Background(), userID, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrAvailabilityOverlap) {
			t.Errorf("expected ErrAvailabilityOverlap, got %v", err)
		}
	})

	t.Run("should fail when new availability completely engulfs existing availability", func(t *testing.T) {
		tc := setUpCreateAvailability(t)

		existingStart, _ := domain.NewTimeOfDay(9, 0)
		existingEnd, _ := domain.NewTimeOfDay(11, 0)
		existingAvail, _ := domain.NewConsultantAvailability(
			"conav_existing",
			"con_123",
			time.Monday,
			existingStart,
			existingEnd,
			tc.clock.Now(),
		)

		tc.availabilityRepo.FindAvailabilitiesByConsultantIDAndDayOfWeekFn = func(ctx context.Context, consultantID string, dayOfWeek time.Weekday) ([]*domain.ConsultantAvailability, error) {
			return []*domain.ConsultantAvailability{existingAvail}, nil
		}

		req := &dto.CreateAvailabilityRequest{
			DayOfWeek: 1,
			StartTime: "08:00",
			EndTime:   "12:00",
		}

		err := tc.sut.Execute(context.Background(), userID, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrAvailabilityOverlap) {
			t.Errorf("expected ErrAvailabilityOverlap, got %v", err)
		}
	})

	t.Run("should fail when ID generator returns error", func(t *testing.T) {
		tc := setUpCreateAvailability(t)
		expectedErr := errors.New("id generator failed")
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

	t.Run("should fail when availability repository SaveAvailability returns error", func(t *testing.T) {
		tc := setUpCreateAvailability(t)
		expectedErr := errors.New("failed to save availability")
		tc.availabilityRepo.SaveAvailabilityFn = func(ctx context.Context, availability *domain.ConsultantAvailability) error {
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
