package outbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

type VerificationRepository interface {
	Save(ctx context.Context, verification *domain.ConsultantVerification) error

	FindByID(
		ctx context.Context,
		id string,
	) (*domain.ConsultantVerification, error)

	FindByConsultantID(
		ctx context.Context,
		consultantID string,
	) (*domain.ConsultantVerification, error)

	Update(ctx context.Context, verification *domain.ConsultantVerification) error
}