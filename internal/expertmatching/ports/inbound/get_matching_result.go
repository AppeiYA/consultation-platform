package inbound

import (
	"context"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
)

type GetMatchingResultRequest struct {
	CaseID string // used to fetch latest run
	RunID  string // optional; used to fetch specific historical run
	TopN   int    // optional; defaults to 5 if <= 0
}

type MatchingResultResponse struct {
	RunID          string
	CaseID         string
	Status         string
	RankingVersion string
	TopCandidates  []domain.RankedCandidate
	TotalRanked    int
	CompletedAt    *time.Time
}

type GetMatchingResultInt interface {
	Execute(ctx context.Context, req GetMatchingResultRequest) (*MatchingResultResponse, error)
}
