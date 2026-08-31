package domain

import custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"

type RankedCandidate struct {
	consultantID string
	rank         Rank
	score        MatchScore
	reasons      []MatchReason
}

var (
	ErrInvalidConsultantID = custom_errors.BadException("invalid consultant ID")
)

func NewRankedCandidate(
	consultantID string,
	rank Rank,
	score MatchScore,
	reasons []MatchReason,
) (RankedCandidate, error) {
	if len(consultantID) == 0 {
		return RankedCandidate{}, ErrInvalidConsultantID
	}

	// Ingress shield: clone the incoming slice
	var clonedReasons []MatchReason
	if len(reasons) > 0 {
		clonedReasons = make([]MatchReason, len(reasons))
		copy(clonedReasons, reasons)
	}

	return RankedCandidate{
		consultantID: consultantID,
		rank:         rank,
		score:        score,
		reasons:      clonedReasons,
	}, nil
}

func (rc RankedCandidate) ConsultantID() string {
	return rc.consultantID
}

func (rc RankedCandidate) Rank() Rank {
	return rc.rank
}

func (rc RankedCandidate) Score() MatchScore {
	return rc.score
}

// Reasons returns a defensive copy so callers cannot mutate internal state
func (rc RankedCandidate) Reasons() []MatchReason {
	if len(rc.reasons) == 0 {
		return nil
	}
	// Egress shield: clone the internal slice
	clonedReasons := make([]MatchReason, len(rc.reasons))
	copy(clonedReasons, rc.reasons)
	return clonedReasons
}
