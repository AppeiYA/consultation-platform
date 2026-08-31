package mocks

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/outbound"
)

type MockCandidateGenerator struct {
	GenerateCandidatesFn func(ctx context.Context, criteria domain.CandidateGenerationCriteria) (domain.CandidatePool, error)
}

func (m *MockCandidateGenerator) GenerateCandidates(ctx context.Context, criteria domain.CandidateGenerationCriteria) (domain.CandidatePool, error) {
	if m.GenerateCandidatesFn != nil {
		return m.GenerateCandidatesFn(ctx, criteria)
	}
	return domain.CandidatePool{}, nil
}

var _ outbound.CandidateGenerator = (*MockCandidateGenerator)(nil)
