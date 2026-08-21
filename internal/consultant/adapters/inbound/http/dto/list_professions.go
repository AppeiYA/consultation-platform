package dto

import "github.com/AppeiYA/consultation-platform/internal/consultant/domain"

type ListProfessionsResponse struct {
	ID  string `json:"id"`
	Name string `json:"name"`	
}

func ProfessionFromDomain(profession *domain.Profession) *ListProfessionsResponse {
	dto := &ListProfessionsResponse{}
	dto.ID = profession.ID()
	dto.Name = profession.Name()
	return dto
}