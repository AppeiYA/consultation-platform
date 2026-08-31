package outbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

type ExpertiseRepository interface {
	SaveMany(ctx context.Context, expertises []*domain.Expertise) error
	Add(ctx context.Context, expertise *domain.Expertise) error
	FindByConsultantID(ctx context.Context, consultantID string) ([]*domain.Expertise, error)
	Delete(ctx context.Context, consultantID string, expertiseID string) error
	ReplaceAll(ctx context.Context, consultantID string, expertises []*domain.Expertise) error
}
