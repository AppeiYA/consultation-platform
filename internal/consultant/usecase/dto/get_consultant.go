package dto

type GetConsultantResponseDto struct {
	ID string
	Profession string
	DisplayName string
	UserID string

	Bio string

	YearsExperience int
	IsAcceptingClients bool
	CreatedAt string
	UpdatedAt string
}