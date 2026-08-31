package outbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
)

type RankingRequest struct {
	CaseDetails    CaseDetails
	RankingVersion domain.RankingVersion
	CandidatePool  domain.CandidatePool
}

type CandidateRanker interface {
	Rank(
		ctx context.Context,
		req RankingRequest,
	) ([]domain.RankedCandidate, error)
}
