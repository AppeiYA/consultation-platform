package domain

import (
	"time"

	custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"
)

var MatchingRunPrefixID = "mrun"

type MatchingRun struct {
	id                 string
	caseID             string
	status             RunStatus
	rankingVersion     RankingVersion
	candidates         []RankedCandidate
	failureReason      string
	cancellationReason string
	createdAt          time.Time
	completedAt        *time.Time
}

var (
	ErrMatchingRunNotFound         = custom_errors.NotFoundError("matching run not found")
	ErrInvalidMatchingRunID       = custom_errors.BadException("invalid matching run ID")
	ErrInvalidCaseID              = custom_errors.BadException("invalid case ID")
	ErrInvalidRunStatusTransition = custom_errors.BadException("invalid matching run status transition")
	ErrRunInTerminalState         = custom_errors.BadException("matching run is in terminal state")
	ErrEmptyFailureReason         = custom_errors.BadException("failure reason cannot be empty")
	ErrEmptyCancellationReason    = custom_errors.BadException("cancellation reason cannot be empty")
	ErrDuplicateCandidate         = custom_errors.BadException("duplicate candidate in matching run")
	ErrDuplicateRank              = custom_errors.BadException("duplicate rank in matching run")
	ErrNonContiguousRank          = custom_errors.BadException("candidate ranks must be contiguous starting from 1")
	ErrContradictoryRankScore     = custom_errors.BadException("higher ranked candidate cannot have a lower match score")
)

// NewMatchingRun initializes a fresh matching run in PENDING state.
func NewMatchingRun(id string, caseID string, rankingVersion RankingVersion, createdAt time.Time) (MatchingRun, error) {
	if len(id) == 0 {
		return MatchingRun{}, ErrInvalidMatchingRunID
	}
	if len(caseID) == 0 {
		return MatchingRun{}, ErrInvalidCaseID
	}

	return MatchingRun{
		id:             id,
		caseID:         caseID,
		status:         RunStatusPending,
		rankingVersion: rankingVersion,
		candidates:     nil,
		createdAt:      createdAt,
		completedAt:    nil,
	}, nil
}

// ReconstituteMatchingRun restores an existing run from persistence.
func ReconstituteMatchingRun(
	id string,
	caseID string,
	status RunStatus,
	rankingVersion RankingVersion,
	candidates []RankedCandidate,
	failureReason string,
	cancellationReason string,
	createdAt time.Time,
	completedAt *time.Time,
) (MatchingRun, error) {
	if len(id) == 0 {
		return MatchingRun{}, ErrInvalidMatchingRunID
	}
	if len(caseID) == 0 {
		return MatchingRun{}, ErrInvalidCaseID
	}
	if !status.IsValid() {
		return MatchingRun{}, ErrInvalidRunStatus
	}

	var clonedCandidates []RankedCandidate
	if len(candidates) > 0 {
		clonedCandidates = make([]RankedCandidate, len(candidates))
		copy(clonedCandidates, candidates)
	}

	return MatchingRun{
		id:                 id,
		caseID:             caseID,
		status:             status,
		rankingVersion:     rankingVersion,
		candidates:         clonedCandidates,
		failureReason:      failureReason,
		cancellationReason: cancellationReason,
		createdAt:          createdAt,
		completedAt:        completedAt,
	}, nil
}

// --- Lifecycle Transitions ---

// StartGeneration transitions run from PENDING -> GENERATING
func (mr *MatchingRun) StartGeneration() error {
	if !mr.status.CanTransitionTo(RunStatusGenerating) {
		return ErrInvalidRunStatusTransition
	}
	mr.status = RunStatusGenerating
	return nil
}

// StartRanking transitions run from GENERATING -> RANKING
func (mr *MatchingRun) StartRanking() error {
	if !mr.status.CanTransitionTo(RunStatusRanking) {
		return ErrInvalidRunStatusTransition
	}
	mr.status = RunStatusRanking
	return nil
}

// Complete transitions run from RANKING -> COMPLETED with validated ranked candidates (0 candidates is valid)
func (mr *MatchingRun) Complete(candidates []RankedCandidate, completedAt time.Time) error {
	if !mr.status.CanTransitionTo(RunStatusCompleted) {
		return ErrInvalidRunStatusTransition
	}

	if len(candidates) > 0 {
		seenConsultants := make(map[string]struct{}, len(candidates))
		for i, c := range candidates {
			// 1. Check duplicate consultant
			if _, exists := seenConsultants[c.ConsultantID()]; exists {
				return ErrDuplicateCandidate
			}
			seenConsultants[c.ConsultantID()] = struct{}{}

			// 2. Check contiguous rank matching position (1-indexed)
			expectedRank := i + 1
			if c.Rank().Value() != expectedRank {
				return ErrNonContiguousRank
			}

			// 3. Check monotonic score (higher rank should have score >= lower rank)
			if i > 0 {
				prevScore := candidates[i-1].Score().Value()
				currScore := c.Score().Value()
				if currScore > prevScore {
					return ErrContradictoryRankScore
				}
			}
		}
	}

	// Defensive copy
	var cloned []RankedCandidate
	if len(candidates) > 0 {
		cloned = make([]RankedCandidate, len(candidates))
		copy(cloned, candidates)
	}

	mr.status = RunStatusCompleted
	mr.candidates = cloned
	mr.completedAt = &completedAt
	return nil
}

// Fail transitions run to FAILED with a required reason
func (mr *MatchingRun) Fail(reason string, completedAt time.Time) error {
	if !mr.status.CanTransitionTo(RunStatusFailed) {
		return ErrInvalidRunStatusTransition
	}
	if len(reason) == 0 {
		return ErrEmptyFailureReason
	}

	mr.status = RunStatusFailed
	mr.failureReason = reason
	mr.completedAt = &completedAt
	return nil
}

// Cancel transitions run to CANCELLED with a required reason
func (mr *MatchingRun) Cancel(reason string, completedAt time.Time) error {
	if !mr.status.CanTransitionTo(RunStatusCancelled) {
		return ErrInvalidRunStatusTransition
	}
	if len(reason) == 0 {
		return ErrEmptyCancellationReason
	}

	mr.status = RunStatusCancelled
	mr.cancellationReason = reason
	mr.completedAt = &completedAt
	return nil
}

// --- Query Methods & Projections ---

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

func (mr MatchingRun) FailureReason() string {
	return mr.failureReason
}

func (mr MatchingRun) CancellationReason() string {
	return mr.cancellationReason
}

func (mr MatchingRun) CreatedAt() time.Time {
	return mr.createdAt
}

func (mr MatchingRun) CompletedAt() *time.Time {
	if mr.completedAt == nil {
		return nil
	}
	t := *mr.completedAt
	return &t
}

func (mr MatchingRun) IsCompleted() bool {
	return mr.status == RunStatusCompleted
}

func (mr MatchingRun) Candidates() []RankedCandidate {
	if len(mr.candidates) == 0 {
		return nil
	}
	cloned := make([]RankedCandidate, len(mr.candidates))
	copy(cloned, mr.candidates)
	return cloned
}

// TopN safely slices the first N ranked candidates (or all if n <= 0 or n >= len)
func (mr MatchingRun) TopN(n int) []RankedCandidate {
	if len(mr.candidates) == 0 {
		return []RankedCandidate{}
	}
	if n <= 0 || n >= len(mr.candidates) {
		cloned := make([]RankedCandidate, len(mr.candidates))
		copy(cloned, mr.candidates)
		return cloned
	}

	cloned := make([]RankedCandidate, n)
	copy(cloned, mr.candidates[:n])
	return cloned
}
