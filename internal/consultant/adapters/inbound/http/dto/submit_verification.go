package dto

type SubmitVerificationResponseDTO struct {
	VerificationID    string `json:"verification_id"`
	ProviderReference string `json:"provider_reference"`
	VerificationURL   string `json:"verification_url"`
	Status            string `json:"status"`
}
