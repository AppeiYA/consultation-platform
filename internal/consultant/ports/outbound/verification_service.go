package outbound

import "context"

type VerificationService interface {
	CreateInquiry(
		ctx context.Context,
		consultantID string,
	) (*CreateInquiryResult, error)
}

type CreateInquiryResult struct {
	Provider string
	ProviderReference string
	VerificationURL   string
}