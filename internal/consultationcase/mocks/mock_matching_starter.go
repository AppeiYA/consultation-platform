package mocks

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultationcase/ports/outbound"
)

type MockMatchingStarter struct {
	StartMatchingFn func(ctx context.Context, caseID string) error
}

func (m *MockMatchingStarter) StartMatching(ctx context.Context, caseID string) error {
	if m.StartMatchingFn != nil {
		return m.StartMatchingFn(ctx, caseID)
	}
	return nil
}

var _ outbound.MatchingStarter = (*MockMatchingStarter)(nil)
