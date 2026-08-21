package mocks

import (
	"context"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

type MockAvailabilityRepository struct {
	SaveAvailabilityFn func(ctx context.Context, availability *domain.ConsultantAvailability) error
	FindAvailabilityByIDFn func(ctx context.Context, id string) (*domain.ConsultantAvailability, error)
	FindAvailabilitiesByConsultantIDFn func(ctx context.Context, consultantID string) ([]*domain.ConsultantAvailability, error)
	FindAvailabilitiesByConsultantIDAndDayOfWeekFn func(ctx context.Context, consultantID string, dayOfWeek time.Weekday) ([]*domain.ConsultantAvailability, error)
	UpdateAvailabilityFn func(ctx context.Context, availability *domain.ConsultantAvailability) error
	DeleteAvailabilityFn func(ctx context.Context, id string) error
}

func (m *MockAvailabilityRepository) SaveAvailability(ctx context.Context, availability *domain.ConsultantAvailability) error {
	if m.SaveAvailabilityFn != nil {
		return m.SaveAvailabilityFn(ctx, availability)
	}
	return nil
}

func (m *MockAvailabilityRepository) FindAvailabilityByID(ctx context.Context, id string) (*domain.ConsultantAvailability, error) {
	if m.FindAvailabilityByIDFn != nil {
		return m.FindAvailabilityByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockAvailabilityRepository) FindAvailabilitiesByConsultantID(ctx context.Context, consultantID string) ([]*domain.ConsultantAvailability, error) {
	if m.FindAvailabilitiesByConsultantIDFn != nil {
		return m.FindAvailabilitiesByConsultantIDFn(ctx, consultantID)
	}
	return nil, nil
}

func (m *MockAvailabilityRepository) FindAvailabilitiesByConsultantIDAndDayOfWeek(ctx context.Context, consultantID string, dayOfWeek time.Weekday) ([]*domain.ConsultantAvailability, error) {
	if m.FindAvailabilitiesByConsultantIDAndDayOfWeekFn != nil {
		return m.FindAvailabilitiesByConsultantIDAndDayOfWeekFn(ctx, consultantID, dayOfWeek)
	}
	return nil, nil
}

func (m *MockAvailabilityRepository) UpdateAvailability(ctx context.Context, availability *domain.ConsultantAvailability) error {
	if m.UpdateAvailabilityFn != nil {
		return m.UpdateAvailabilityFn(ctx, availability)
	}
	return nil
}

func (m *MockAvailabilityRepository) DeleteAvailability(ctx context.Context, id string) error {
	if m.DeleteAvailabilityFn != nil {
		return m.DeleteAvailabilityFn(ctx, id)
	}
	return nil
}
