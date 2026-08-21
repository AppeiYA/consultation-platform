package consultantRepo

import (
	"time"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

type Consultant struct {
	ID   string `db:"id"`
	UserID string `db:"user_id"`
	ProfessionID string `db:"profession_id"`
	DisplayName string `db:"display_name"`
	Bio string `db:"bio"`
	YearsExperience int `db:"years_experience"`
	IsAcceptingClients bool `db:"is_accepting_clients"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (c *Consultant) ToDomain() (*domain.Consultant, error) {
	return domain.ReconstitueConsultant(
		c.ID,
		c.UserID,
		c.ProfessionID,
		c.DisplayName,
		c.Bio,
		c.YearsExperience,
		c.IsAcceptingClients,
		c.CreatedAt,
		c.UpdatedAt,
	)
}

func ConsultantFromDomain(consultant *domain.Consultant) *Consultant {
	return &Consultant{
		ID: consultant.ID(),
		UserID: consultant.UserID(),
		ProfessionID: string(consultant.ProfessionID()),
		DisplayName: consultant.DisplayName().String(),
		Bio: consultant.Bio().String(),
		YearsExperience: consultant.YearsExperience().Int(),
		IsAcceptingClients: consultant.IsAcceptingClients(),
		CreatedAt: consultant.CreatedAt(),
		UpdatedAt: consultant.UpdatedAt(),
	}
}