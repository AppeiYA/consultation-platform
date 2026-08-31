package inbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
)

type StartMatchingInt interface {
	Execute(ctx context.Context, caseID string) (*domain.MatchingRun, error)
}
