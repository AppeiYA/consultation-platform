package expertMatchingRepo

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
)

type MatchingRunModel struct {
	ID                 string         `db:"id"`
	CaseID             string         `db:"case_id"`
	Status             string         `db:"status"`
	RankingVersion     string         `db:"ranking_version"`
	FailureReason      sql.NullString `db:"failure_reason"`
	CancellationReason sql.NullString `db:"cancellation_reason"`
	CreatedAt          time.Time      `db:"created_at"`
	CompletedAt        *time.Time     `db:"completed_at"`
}

type CandidateModel struct {
	ID           int     `db:"id"`
	RunID        string  `db:"run_id"`
	ConsultantID string  `db:"consultant_id"`
	RankPosition int     `db:"rank_position"`
	Score        float64 `db:"score"`
	ReasonsJSON  string  `db:"reasons"`
}

type ReasonJSON struct {
	Factor string `json:"factor"`
	Detail string `json:"detail"`
}

func (m *MatchingRunModel) ToDomain(candidateModels []CandidateModel) (*domain.MatchingRun, error) {
	status, err := domain.NewRunStatusFromString(m.Status)
	if err != nil {
		return nil, err
	}

	rankingVersion, err := domain.NewRankingVersion(m.RankingVersion)
	if err != nil {
		return nil, err
	}

	candidates := make([]domain.RankedCandidate, 0, len(candidateModels))
	for _, cm := range candidateModels {
		rank, err := domain.NewRank(cm.RankPosition)
		if err != nil {
			return nil, err
		}
		score, err := domain.NewMatchScore(cm.Score)
		if err != nil {
			return nil, err
		}

		var reasonJSONs []ReasonJSON
		if len(cm.ReasonsJSON) > 0 && cm.ReasonsJSON != "[]" {
			if err := json.Unmarshal([]byte(cm.ReasonsJSON), &reasonJSONs); err != nil {
				return nil, err
			}
		}

		reasons := make([]domain.MatchReason, 0, len(reasonJSONs))
		for _, rj := range reasonJSONs {
			factor, err := domain.NewMatchFactorFromString(rj.Factor)
			if err != nil {
				return nil, err
			}
			reason, err := domain.NewMatchReason(factor, rj.Detail)
			if err != nil {
				return nil, err
			}
			reasons = append(reasons, reason)
		}

		cand, err := domain.NewRankedCandidate(cm.ConsultantID, rank, score, reasons)
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, cand)
	}

	failureReason := ""
	if m.FailureReason.Valid {
		failureReason = m.FailureReason.String
	}
	cancellationReason := ""
	if m.CancellationReason.Valid {
		cancellationReason = m.CancellationReason.String
	}

	run, err := domain.ReconstituteMatchingRun(
		m.ID,
		m.CaseID,
		status,
		rankingVersion,
		candidates,
		failureReason,
		cancellationReason,
		m.CreatedAt,
		m.CompletedAt,
	)
	if err != nil {
		return nil, err
	}

	return &run, nil
}
