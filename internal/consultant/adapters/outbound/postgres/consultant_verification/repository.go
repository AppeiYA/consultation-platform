package consultantVerificationRepo

import (
	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
	system "github.com/AppeiYA/consultation-platform/internal/shared/adapters/outbound/clock"
	"github.com/AppeiYA/consultation-platform/internal/shared/db"
)


type VerificationRepository struct {
	repository db.Repository
	clock *system.SystemClock
}

func NewVerificationRepository(repository db.Repository, clock *system.SystemClock) *VerificationRepository {
	return &VerificationRepository{
		repository: repository,
		clock: clock,
	}
}

var _ outbound.VerificationRepository = (*VerificationRepository)(nil)