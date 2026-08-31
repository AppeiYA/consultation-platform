package expertMatchingRepo

import (
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/shared/db"
)

type ExpertMatchingRepository struct {
	repository *db.Repository
}

func NewExpertMatchingRepository(repository *db.Repository) *ExpertMatchingRepository {
	return &ExpertMatchingRepository{
		repository: repository,
	}
}

var _ outbound.MatchingRunRepository = (*ExpertMatchingRepository)(nil)