package dto

import (
	"fmt"
	usecase_dto "github.com/AppeiYA/consultation-platform/internal/consultationcase/usecase/dto"
)

type CreateConsultationCaseDTO struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Category     string `json:"category"`
}

func (dto *CreateConsultationCaseDTO) Validate() error {
	if dto.Title == "" {
		return fmt.Errorf("title is required")
	}
	if dto.Description == "" {
		return fmt.Errorf("description is required")
	}
	if dto.Category == "" {
		return fmt.Errorf("category is required")
	}
	return nil
}

func (dto *CreateConsultationCaseDTO) ToUsecaseDTO() *usecase_dto.CreateCaseDTO {
	return &usecase_dto.CreateCaseDTO{
		Title:       dto.Title,
		Description: dto.Description,
		Category:    dto.Category,
	}
}