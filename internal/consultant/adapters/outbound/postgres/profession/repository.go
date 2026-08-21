package professionRepo

import (
	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
	system "github.com/AppeiYA/consultation-platform/internal/shared/adapters/outbound/clock"
	"github.com/AppeiYA/consultation-platform/internal/shared/db"
)

type ProfessionRepository struct {
	repository db.Repository
	clock *system.SystemClock
}

func NewProfessionRepository(repository db.Repository, clock *system.SystemClock) *ProfessionRepository {
	return &ProfessionRepository{
		repository: repository,
		clock: clock,
	}
}

var _ outbound.ProfessionRepository = (*ProfessionRepository)(nil)