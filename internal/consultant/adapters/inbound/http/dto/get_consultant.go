package dto

import (
	"github.com/AppeiYA/consultation-platform/internal/consultant/usecase/dto"
)

type PublicConsultantResponseDTO struct {
	ID   string `json:"id"`
	Profession string `json:"profession"`
	DisplayName string `json:"display_name"`
	UserID string `json:"user_id"`

	Bio string `json:"bio"`

	YearsExperience int `json:"years_experience"`
	IsAcceptingClients bool `json:"is_accepting_clients"`
	CreatedAt string `json:"created_at"`
}

type PrivateConsultantResponseDTO struct {
	PublicConsultantResponseDTO

	UpdatedAt string `json:"updated_at"`
}

func (dto *PublicConsultantResponseDTO) FromUsecaseDTO(consultant *dto.GetConsultantResponseDto) {
	dto.ID = consultant.ID
	dto.Profession = consultant.Profession
	dto.DisplayName = consultant.DisplayName
	dto.UserID = consultant.UserID
	dto.Bio = consultant.Bio
	dto.YearsExperience = consultant.YearsExperience
	dto.IsAcceptingClients = consultant.IsAcceptingClients
	dto.CreatedAt = consultant.CreatedAt
}

func (dto *PrivateConsultantResponseDTO) FromUsecaseDTO(consultant *dto.GetConsultantResponseDto) {
	dto.PublicConsultantResponseDTO.FromUsecaseDTO(consultant)
	dto.UpdatedAt = consultant.UpdatedAt
}