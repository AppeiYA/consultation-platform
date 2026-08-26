package domain

import (
	"time"

	custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"
)

type MatchingRun struct {
    id               string
    caseID           string
    status           RunStatus
    rankingVersion   RankingVersion
    candidates       []RankedCandidate
    createdAt        time.Time
    completedAt      *time.Time
}

var (
	ErrInvalidMatchingRunID = custom_errors.BadException("invalid matching run ID")
	ErrInvalidCaseID = custom_errors.BadException("invalid case ID")
	ErrInvalidCompletedAt = custom_errors.BadException("completedAt must be nil for a new matching run")
	ErrEmptyCandidates = custom_errors.BadException("candidates list cannot be empty")
)
func NewMatchingRun(id string, caseID string, status RunStatus, rankingVersion RankingVersion, candidates []RankedCandidate, createdAt time.Time, completedAt *time.Time) (MatchingRun, error) {
	if id == "" {
		return MatchingRun{}, ErrInvalidMatchingRunID
	}
	if caseID == "" {
		return MatchingRun{}, ErrInvalidCaseID
	}
	if status != RunStatusPending {
		return MatchingRun{}, ErrInvalidRunStatus
	}

	if completedAt != nil { 
		return MatchingRun{}, ErrInvalidCompletedAt
	}

	if len(candidates) == 0 {
		return MatchingRun{}, ErrEmptyCandidates
	}
	return MatchingRun{
		id:             id,
		caseID:         caseID,
		status:         status,
		rankingVersion: rankingVersion,
		candidates:     candidates,
		createdAt:      createdAt,
		completedAt:    completedAt,
	}, nil
}

func (mr MatchingRun) ID() string {
	return mr.id
}

func (mr MatchingRun) CaseID() string {
	return mr.caseID
}

func (mr MatchingRun) Status() RunStatus {
	return mr.status
}

func (mr MatchingRun) RankingVersion() RankingVersion {
	return mr.rankingVersion
}

func (mr MatchingRun) Candidates() []RankedCandidate {
	if len(mr.candidates) == 0 {
		return nil
	}
	clonedCandidates := make([]RankedCandidate, len(mr.candidates))
	copy(clonedCandidates, mr.candidates)
	return clonedCandidates
}

func (mr MatchingRun) CreatedAt() time.Time {
	return mr.createdAt
}

func (mr MatchingRun) CompletedAt() *time.Time {
	if mr.completedAt == nil {
		return nil
	}
	completedAtCopy := *mr.completedAt
	return &completedAtCopy
}