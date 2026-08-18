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
	IsAcceptingClients bool `db:"is_accepting_clients"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type ConsultantVerification struct {
	ID string `db:"id"`
	ConsultantID string `db:"consultant_id"`
	Provider string `db:"provider"`
	ProviderReference string `db:"provider_reference"`
	Status string `db:"status"`
	SubmittedAt *time.Time `db:"submitted_at"`
	CompletedAt *time.Time `db:"completed_at"`
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
		IsAcceptingClients: consultant.IsAcceptingClients(),
		CreatedAt: consultant.CreatedAt(),
		UpdatedAt: consultant.UpdatedAt(),
	}
}

func (v *ConsultantVerification) ToDomain() (*domain.ConsultantVerification, error) {
	status, err := domain.NewVerificationStatus(v.Status)
	if err != nil {
		return nil, err
	}

	return domain.ReconstitueConsultantVerification(
		v.ID,
		v.ConsultantID,
		v.Provider,
		v.ProviderReference,
		status,
		v.SubmittedAt,
		v.CompletedAt,
		v.CreatedAt,
		v.UpdatedAt,
	)
}

func ConsultantVerificationFromDomain(verification *domain.ConsultantVerification) *ConsultantVerification {
	return &ConsultantVerification{
		ID: verification.ID(),
		ConsultantID: verification.ConsultantID(),
		Provider: verification.Provider(),
		ProviderReference: verification.ProviderReference(),
		Status: string(verification.Status()),
		SubmittedAt: verification.SubmittedAt(),
		CompletedAt: verification.CompletedAt(),
		CreatedAt: verification.CreatedAt(),
		UpdatedAt: verification.UpdatedAt(),
	}
}