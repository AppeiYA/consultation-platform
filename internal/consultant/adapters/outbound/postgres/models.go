package postgres

import (
	"time"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

type Consultant struct {
	ID   string `db:"id"`
	UserID string `db:"user_id"`
	Profession string `db:"profession"`
	DisplayName string `db:"display_name"`
	Bio string `db:"bio"`
	YearsExperience int `db:"years_experience"`
	IsVerified bool `db:"is_verified"`
	IsAcceptingClients bool `db:"is_accepting_clients"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (c *Consultant) ToDomain() (*domain.Consultant, error) {
	return domain.ReconstitueConsultant(
		c.ID,
		c.UserID,
		c.Profession,
		c.DisplayName,
		c.Bio,
		c.YearsExperience,
		c.IsVerified,
		c.IsAcceptingClients,
		c.CreatedAt,
		c.UpdatedAt,
	)
}

func ConsultantFromDomain(consultant *domain.Consultant) *Consultant {
	return &Consultant{
		ID: consultant.ID(),
		UserID: consultant.UserID(),
		Profession: string(consultant.Profession()),
		DisplayName: consultant.DisplayName().String(),
		Bio: consultant.Bio().String(),
		YearsExperience: consultant.YearsExperience().Int(),
		IsVerified: consultant.IsVerified(),
		IsAcceptingClients: consultant.IsAcceptingClients(),
		CreatedAt: consultant.CreatedAt(),
		UpdatedAt: consultant.UpdatedAt(),
	}
}