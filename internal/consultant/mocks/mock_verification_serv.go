package mocks

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
)

type MockVerificationService struct {
	CreateInquiryFn func(ctx context.Context, consultantID string) (*outbound.CreateInquiryResult, error)
}

func (m *MockVerificationService) CreateInquiry(
	ctx context.Context,
	consultantID string,
) (*outbound.CreateInquiryResult, error) {
	if m.CreateInquiryFn != nil {
		return m.CreateInquiryFn(ctx, consultantID)
	}
	return nil, nil
}