package dto

import usecase_dto "github.com/AppeiYA/consultation-platform/internal/consultant/usecase/dto"

type RegisterConsultantDTO struct {
	Profession string `json:"profession"`
	DisplayName string `json:"display_name"`

	Bio string `json:"bio"`

	YearsExperience int `json:"years_experience"`
}

func (dto *RegisterConsultantDTO) ToUsecaseDTO() *usecase_dto.RegisterConsultantDTO {
	return &usecase_dto.RegisterConsultantDTO{
		Profession: dto.Profession,
		DisplayName: dto.DisplayName,
		Bio: dto.Bio,
		YearsExperience: dto.YearsExperience,
	}
}