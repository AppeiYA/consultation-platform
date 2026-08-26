package consultationCaseRepo

import (
	"github.com/AppeiYA/consultation-platform/internal/consultationcase/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/shared/db"
)

type ConsultationCaseRepository struct {
	repository *db.Repository
}

func NewConsultationCaseRepository(repository *db.Repository) *ConsultationCaseRepository {
	return &ConsultationCaseRepository{
		repository: repository,
	}
}

var _ outbound.CaseRepository = (*ConsultationCaseRepository)(nil)