package mocks

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/outbound"
)

type MockCandidateRanker struct {
	RankFn func(ctx context.Context, req outbound.RankingRequest) ([]domain.RankedCandidate, error)
}

func (m *MockCandidateRanker) Rank(ctx context.Context, req outbound.RankingRequest) ([]domain.RankedCandidate, error) {
	if m.RankFn != nil {
		return m.RankFn(ctx, req)
	}
	return nil, nil
}

var _ outbound.CandidateRanker = (*MockCandidateRanker)(nil)
