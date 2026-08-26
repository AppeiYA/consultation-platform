package domain

import custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"

type RankedCandidate struct {
	consultantID string
	score MatchScore
	reasons []MatchReason
}

var (
	ErrInvalidConsultantID = custom_errors.BadException("invalid consultant ID")
)
func NewRankedCandidate(consultantID string, score MatchScore, reasons []MatchReason) (RankedCandidate, error) {
	if len(consultantID) == 0 {
		return  RankedCandidate{}, ErrInvalidConsultantID
	}
	// Ingress shield
	var clonedReasons []MatchReason
	if len(reasons) > 0 {
		clonedReasons = make([]MatchReason, len(reasons))
		copy(clonedReasons, reasons)
	}

	return RankedCandidate{
		consultantID: consultantID,
		score: score,
		reasons: clonedReasons,
	}, nil
}

func (rc RankedCandidate) ConsultantID() string {
	return rc.consultantID
}

func (rc RankedCandidate) Score() MatchScore {
	return rc.score
}

func (rc RankedCandidate) Reasons() []MatchReason {
	if len(rc.reasons) == 0 {
		return nil
	}
	// Egress shield
	clonedReasons := make([]MatchReason, len(rc.reasons))
	copy(clonedReasons, rc.reasons)
	return clonedReasons
}