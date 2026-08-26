package outbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultationcase/domain"
)

type CaseRepository interface {
	SaveCase(ctx context.Context, c *domain.ConsultationCase) error
	FindCaseByID(ctx context.Context, id string) (*domain.ConsultationCase, error)
	FindCasesByClientID(ctx context.Context, clientID string) ([]*domain.ConsultationCase, error)
	UpdateCase(ctx context.Context, c *domain.ConsultationCase) error
	DeleteCase(ctx context.Context, id string) error
}