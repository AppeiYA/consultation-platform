package consultationCaseRepo

import (
	"time"

	"github.com/AppeiYA/consultation-platform/internal/consultationcase/domain"
)

type CaseModel struct {
	ID string `db:"id"`
	ClientID string `db:"client_id"`
	Title string `db:"title"`
	Description string `db:"description"`
	Category string `db:"category"`
	Status string `db:"status"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (c *CaseModel) ToDomain() (*domain.ConsultationCase, error) {
	return domain.ReconstituteConsultationCase(
		c.ID,
		c.ClientID,
		c.Title,
		c.Description,
		c.Category,
		c.Status, 
		c.CreatedAt,
		c.UpdatedAt,
	)
}

func FromDomainToModel(consultationCase *domain.ConsultationCase) *CaseModel {
	return &CaseModel{
		ID: consultationCase.ID(),
		ClientID: consultationCase.ClientID(),
		Title: consultationCase.Title().String(),
		Description: consultationCase.Description().String(),
		Category: consultationCase.Category().String(),
		Status: string(consultationCase.Status()),
		CreatedAt: consultationCase.CreatedAt(),
		UpdatedAt: consultationCase.UpdatedAt(),
	}
}