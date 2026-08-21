package professionRepo

import (
	"time"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

type ProfessionModel struct {
	ID string `db:"id"`
	Name string `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

func (p *ProfessionModel) ToDomain() domain.Profession {
	return domain.ReconstituteProfession(p.ID, p.Name, p.CreatedAt)
}

func FromDomain(profession domain.Profession) *ProfessionModel {
	return &ProfessionModel{
		ID: profession.ID(),
		Name: profession.Name(),
		CreatedAt: profession.CreatedAt(),
	}
}