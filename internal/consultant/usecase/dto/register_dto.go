package dto

type RegisterConsultantDTO struct {
	ProfessionID    string
	DisplayName     string
	Bio             string
	YearsExperience int
	Expertises      []string
}
