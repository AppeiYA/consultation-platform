package domain

import custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"

type RunStatus string

var (
	ErrInvalidRunStatus = custom_errors.BadException("invalid run status")
)
const (
	RunStatusPending   RunStatus = "pending"
	RunStatusGenerating RunStatus = "generating"
	RunStatusRanking	RunStatus = "ranking"
	RunStatusCompleted RunStatus = "completed"
	RunStatusFailed    RunStatus = "failed"
	RunStatusCancelled RunStatus = "cancelled"
)

var validTransitions = map[RunStatus][]RunStatus{
	RunStatusPending:   {RunStatusGenerating, RunStatusCancelled},
	RunStatusGenerating: {RunStatusRanking, RunStatusFailed, RunStatusCancelled},
	RunStatusRanking:   {RunStatusCompleted, RunStatusFailed, RunStatusCancelled},
	RunStatusCompleted: {},
	RunStatusFailed:    {RunStatusGenerating},
	RunStatusCancelled: {},
}

func (rs RunStatus) CanTransitionTo(newStatus RunStatus) bool {
	validNextStatuses, exists := validTransitions[rs]
	if !exists {
		return false
	}
	for _, status := range validNextStatuses {
		if status == newStatus {
			return true
		}
	}
	return false
}

func (rs RunStatus) IsValid() bool {
	switch rs {
	case RunStatusPending, RunStatusGenerating, RunStatusRanking, RunStatusCompleted, RunStatusFailed, RunStatusCancelled:
		return true
	default:
		return false
	}
}

func NewRunStatusFromString(s string) (RunStatus, error) {
	rs := RunStatus(s)
	if !rs.IsValid() {
		return "", ErrInvalidRunStatus
	}
	return rs, nil
}