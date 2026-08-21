package consultantRepo

import (
	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
	system "github.com/AppeiYA/consultation-platform/internal/shared/adapters/outbound/clock"
	"github.com/AppeiYA/consultation-platform/internal/shared/db"
)

type ConsultantRepository struct {
	repository db.Repository
	clock *system.SystemClock
}

func NewConsultantRepository(repository db.Repository, clock *system.SystemClock) *ConsultantRepository {
	return &ConsultantRepository{
		repository: repository,
		clock: clock,
	}
}

var _ outbound.ConsultantRepository = (*ConsultantRepository)(nil)