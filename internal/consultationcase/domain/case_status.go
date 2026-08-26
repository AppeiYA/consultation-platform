package domain

import custom_errors "github.com/AppeiYA/consultation-platform/internal/shared/errors"

type CaseStatus string 
const (
	CaseStatusDraft       CaseStatus = "DRAFT"
	CaseStatusSubmitted   CaseStatus = "SUBMITTED"
	CaseStatusMatching    CaseStatus = "MATCHING"
	CaseStatusMatched     CaseStatus = "MATCHED"
	CaseStatusInProgress  CaseStatus = "IN_PROGRESS"
	CaseStatusResolved    CaseStatus = "RESOLVED"
	CaseStatusCancelled   CaseStatus = "CANCELLED"
)

var (
	ErrCaseStatusInvalid = custom_errors.BadException("Invalid case status")
)

var statusTransitionMap = map[CaseStatus][]CaseStatus {
	CaseStatusDraft:       {CaseStatusSubmitted},
	CaseStatusSubmitted:   {CaseStatusMatching},
	CaseStatusMatching:    {CaseStatusMatched, CaseStatusCancelled},
	CaseStatusMatched:     {CaseStatusInProgress},
	CaseStatusInProgress:  {CaseStatusResolved, CaseStatusCancelled},
	CaseStatusResolved:    {},
	CaseStatusCancelled:   {},
}

func NewCaseStatus(value string) (CaseStatus, error) {
	cs := CaseStatus(value)
	if !cs.isValid() {
		return "", ErrCaseStatusInvalid
	}
	return cs, nil
}

func (cs CaseStatus) isValid() bool {
	switch cs {
	case CaseStatusDraft, CaseStatusSubmitted, CaseStatusMatching, CaseStatusMatched, CaseStatusInProgress, CaseStatusResolved, CaseStatusCancelled:
		return true
	default:
		return false
	}
}

func (cs CaseStatus) CanTransitionTo(newStatus CaseStatus) bool {
	validTransitions, exists := statusTransitionMap[cs]
	if !exists {
		return false
	}
	for _, validStatus := range validTransitions {
		if validStatus == newStatus {
			return true
		}
	}
	return false
}