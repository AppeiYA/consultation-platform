package dto

type SubmitVerificationRequest struct {
	ConsultantID string
}

type SubmitVerificationResponse struct {
	VerificationID   string
	ProviderReference string
	VerificationURL   string
	Status            string
}