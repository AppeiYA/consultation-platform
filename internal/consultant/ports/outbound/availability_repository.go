package outbound

import (
	"context"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

type AvailabilityRepository interface {
	SaveAvailability(
		ctx context.Context,
		availability *domain.ConsultantAvailability,
	) error

	FindAvailabilityByID(
		ctx context.Context,
		id string,
	) (*domain.ConsultantAvailability, error)

	FindAvailabilitiesByConsultantID(
		ctx context.Context,
		consultantID string,
	) ([]*domain.ConsultantAvailability, error)

	FindAvailabilitiesByConsultantIDAndDayOfWeek(
		ctx context.Context,
		consultantID string,
		dayOfWeek time.Weekday,
	) ([]*domain.ConsultantAvailability, error)

	UpdateAvailability(
		ctx context.Context,
		availability *domain.ConsultantAvailability,
	) error

	DeleteAvailability(
		ctx context.Context,
		id string,
	) error

	DeactivateAvailability(
		ctx context.Context,
		id string,
	) error

	ActivateAvailability(
		ctx context.Context,
		id string,
	) error
}