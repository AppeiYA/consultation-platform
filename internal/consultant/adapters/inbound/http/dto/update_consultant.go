package dto

import "github.com/AppeiYA/consultation-platform/internal/consultant/usecase/dto"

type UpdateConsultantModel struct {
	ProfessionID        string 	`json:"profession_id"`
	DisplayName       string 	`json:"display_name"`
	Bio               string 	`json:"bio"`
	YearsExperience   int    	`json:"years_experience"`
}

func (u *UpdateConsultantModel) ToUsecaseDTO() dto.UpdateConsultantDTO {
	return dto.UpdateConsultantDTO{
		ProfessionID:      u.ProfessionID,
		DisplayName:     u.DisplayName,
		Bio:             u.Bio,
		YearsExperience: u.YearsExperience,
	}
}