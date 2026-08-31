package mocks

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/outbound"
)

type MockMatchingJobEnqueuer struct {
	EnqueueFn func(ctx context.Context, job outbound.MatchingJob) error
}

func (m *MockMatchingJobEnqueuer) Enqueue(ctx context.Context, job outbound.MatchingJob) error {
	if m.EnqueueFn != nil {
		return m.EnqueueFn(ctx, job)
	}
	return nil
}

var _ outbound.MatchingJobEnqueuer = (*MockMatchingJobEnqueuer)(nil)
