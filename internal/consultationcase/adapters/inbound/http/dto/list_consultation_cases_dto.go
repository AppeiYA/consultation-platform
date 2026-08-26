package dto

import (
	"time"

	"github.com/AppeiYA/consultation-platform/internal/consultationcase/domain"
)

type ConsultationCasesDTO struct {
	ID string `json:"id"`
	ClientID string `json:"client_id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Category string `json:"category"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FromDomainToConsultationCase(cases *domain.ConsultationCase) ConsultationCasesDTO {
	return ConsultationCasesDTO{
		ID: cases.ID(),
		ClientID: cases.ClientID(),
		Title: cases.Title().String(),
		Description: cases.Description().String(),
		Category: cases.Category().String(),
		Status: string(cases.Status()),
		CreatedAt: cases.CreatedAt(),
		UpdatedAt: cases.UpdatedAt(),
	}
}