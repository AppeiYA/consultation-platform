package dto

type AddExpertiseDTO struct {
	Name string `json:"name" validate:"required"`
}

type ReplaceExpertisesDTO struct {
	Expertises []string `json:"expertises" validate:"required"`
}

type ExpertiseResponseDTO struct {
	ID           string `json:"id"`
	ConsultantID string `json:"consultant_id"`
	Name         string `json:"name"`
}
