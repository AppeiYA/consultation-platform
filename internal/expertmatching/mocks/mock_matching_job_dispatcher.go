package mocks

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/outbound"
)

type MockMatchingJobDispatcher struct {
	DispatchMatchingFn func(ctx context.Context, runID string) error
}

func (m *MockMatchingJobDispatcher) DispatchMatching(ctx context.Context, runID string) error {
	if m.DispatchMatchingFn != nil {
		return m.DispatchMatchingFn(ctx, runID)
	}
	return nil
}

var _ outbound.MatchingJobDispatcher = (*MockMatchingJobDispatcher)(nil)
