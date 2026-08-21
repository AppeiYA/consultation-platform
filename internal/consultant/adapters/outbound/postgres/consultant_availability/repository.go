package consultantAvailabilityRepo

import (
	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
	system "github.com/AppeiYA/consultation-platform/internal/shared/adapters/outbound/clock"
	"github.com/AppeiYA/consultation-platform/internal/shared/db"
)

type AvailabilityRepository struct {
	repository db.Repository
	clock *system.SystemClock
}

func NewAvailabilityRepository(repository db.Repository, clock *system.SystemClock) *AvailabilityRepository {
	return &AvailabilityRepository{
		repository: repository,
		clock: clock,
	}
}

var _ outbound.AvailabilityRepository = (*AvailabilityRepository)(nil)