package dto

type ExpertiseResponseDTO struct {
	ID           string `json:"id"`
	ConsultantID string `json:"consultant_id"`
	Name         string `json:"name"`
}

type AddExpertiseDTO struct {
	Name string `json:"name"`
}

type ReplaceExpertisesDTO struct {
	Expertises []string `json:"expertises"`
}
