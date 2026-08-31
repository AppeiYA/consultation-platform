package outbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
)

type MatchingRunRepository interface {
	Save(ctx context.Context, run *domain.MatchingRun) error
	FindByID(ctx context.Context, id string) (*domain.MatchingRun, error)
	FindLatestByCaseID(ctx context.Context, caseID string) (*domain.MatchingRun, error)
	FindByCaseID(ctx context.Context, caseID string) ([]*domain.MatchingRun, error)
}
