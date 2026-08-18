package verification

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
)

type UnavailableVerificationService struct{}

func (s *UnavailableVerificationService) CreateInquiry(
    ctx context.Context,
    consultantID string,
) (*outbound.CreateInquiryResult, error) {
    return nil, domain.ErrVerificationUnavailable
}