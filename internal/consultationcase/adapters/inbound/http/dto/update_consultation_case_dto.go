package dto

import (
	"fmt"

	usecase_dto "github.com/AppeiYA/consultation-platform/internal/consultationcase/usecase/dto"
)

type UpdateConsultationCaseDTO struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Category    *string `json:"category,omitempty"`
}

func (dto *UpdateConsultationCaseDTO) Validate() error {
	if dto.Title == nil && dto.Description == nil && dto.Category == nil {
		return fmt.Errorf("at least one field must be provided for update")
	}
	if dto.Title != nil && *dto.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if dto.Description != nil && *dto.Description == "" {
		return fmt.Errorf("description cannot be empty")
	}
	if dto.Category != nil && *dto.Category == "" {
		return fmt.Errorf("category cannot be empty")
	}
	return nil
}

func (dto *UpdateConsultationCaseDTO) ToUsecaseDTO() *usecase_dto.UpdateCaseDTO {
	return &usecase_dto.UpdateCaseDTO{
		Title:       dto.Title,
		Description: dto.Description,
		Category:    dto.Category,
	}
}
