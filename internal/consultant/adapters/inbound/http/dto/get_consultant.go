package dto

import (
	"time"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

type PublicConsultantResponseDTO struct {
	ID   string `json:"id"`
	Profession string `json:"profession"`
	DisplayName string `json:"display_name"`
	UserID string `json:"user_id"`

	Bio string `json:"bio"`

	YearsExperience int `json:"years_experience"`
	IsVerified bool `json:"is_verified"`
	IsAcceptingClients bool `json:"is_accepting_clients"`
	CreatedAt string `json:"created_at"`
}

type PrivateConsultantResponseDTO struct {
	PublicConsultantResponseDTO

	UpdatedAt string `json:"updated_at"`
}

func (dto *PublicConsultantResponseDTO) FromDomain(consultant *domain.Consultant) {
	dto.ID = consultant.ID()
	dto.UserID = consultant.UserID()
	dto.Profession = string(consultant.Profession())
	dto.DisplayName = consultant.DisplayName().String()
	dto.Bio = consultant.Bio().String()
	dto.YearsExperience = consultant.YearsExperience().Int()
	dto.IsVerified = consultant.IsVerified()
	dto.IsAcceptingClients = consultant.IsAcceptingClients()
	dto.CreatedAt = consultant.CreatedAt().Format(time.RFC3339)
}

func (dto *PrivateConsultantResponseDTO) FromDomain(consultant *domain.Consultant) {
	dto.PublicConsultantResponseDTO.FromDomain(consultant)
	dto.UpdatedAt = consultant.UpdatedAt().Format(time.RFC3339)
}