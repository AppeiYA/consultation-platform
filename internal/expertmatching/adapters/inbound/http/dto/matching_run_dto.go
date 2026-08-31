package dto

import (
	"time"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/inbound"
)

type StartMatchingResponse struct {
	RunID     string    `json:"run_id"`
	CaseID    string    `json:"case_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func FromDomainToStartMatchingResponse(run *domain.MatchingRun) StartMatchingResponse {
	return StartMatchingResponse{
		RunID:     run.ID(),
		CaseID:    run.CaseID(),
		Status:    string(run.Status()),
		CreatedAt: run.CreatedAt(),
	}
}

type MatchReasonDTO struct {
	Factor string `json:"factor"`
	Detail string `json:"detail"`
}

type RankedCandidateDTO struct {
	ConsultantID string           `json:"consultant_id"`
	Rank         int              `json:"rank"`
	Score        float64          `json:"score"`
	Reasons      []MatchReasonDTO `json:"reasons"`
}

type MatchingResultResponse struct {
	RunID          string               `json:"run_id"`
	CaseID         string               `json:"case_id"`
	Status         string               `json:"status"`
	RankingVersion string               `json:"ranking_version"`
	TopCandidates  []RankedCandidateDTO `json:"top_candidates"`
	TotalRanked    int                  `json:"total_ranked"`
	CompletedAt    *time.Time           `json:"completed_at,omitempty"`
}

func FromInboundToMatchingResultResponse(resp *inbound.MatchingResultResponse) MatchingResultResponse {
	candidates := make([]RankedCandidateDTO, 0, len(resp.TopCandidates))
	for _, c := range resp.TopCandidates {
		reasons := make([]MatchReasonDTO, 0, len(c.Reasons()))
		for _, r := range c.Reasons() {
			reasons = append(reasons, MatchReasonDTO{
				Factor: r.Factor().String(),
				Detail: r.Detail(),
			})
		}
		candidates = append(candidates, RankedCandidateDTO{
			ConsultantID: c.ConsultantID(),
			Rank:         c.Rank().Value(),
			Score:        c.Score().Value(),
			Reasons:      reasons,
		})
	}

	return MatchingResultResponse{
		RunID:          resp.RunID,
		CaseID:         resp.CaseID,
		Status:         resp.Status,
		RankingVersion: resp.RankingVersion,
		TopCandidates:  candidates,
		TotalRanked:    resp.TotalRanked,
		CompletedAt:    resp.CompletedAt,
	}
}
