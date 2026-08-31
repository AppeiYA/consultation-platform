package inbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
)

type ExecuteMatchingInt interface {
	Execute(ctx context.Context, runID string) (*domain.MatchingRun, error)
}
