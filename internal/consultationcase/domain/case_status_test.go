package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCaseStatus(t *testing.T) {
	validStatuses := []string{
		"DRAFT",
		"SUBMITTED",
		"MATCHING",
		"MATCHED",
		"IN_PROGRESS",
		"RESOLVED",
		"CANCELLED",
	}

	for _, s := range validStatuses {
		t.Run("valid status "+s, func(t *testing.T) {
			status, err := NewCaseStatus(s)
			require.NoError(t, err)
			require.Equal(t, CaseStatus(s), status)
		})
	}

	t.Run("invalid status", func(t *testing.T) {
		status, err := NewCaseStatus("INVALID_STATUS")
		require.Error(t, err)
		require.Empty(t, status)
		require.Equal(t, ErrCaseStatusInvalid, err)
	})
}

func TestStatusTransitions(t *testing.T) {
	tests := []struct {
		name      string
		current   CaseStatus
		next      CaseStatus
		expectCan bool
	}{
		{"DRAFT to SUBMITTED", CaseStatusDraft, CaseStatusSubmitted, true},
		{"DRAFT to MATCHING", CaseStatusDraft, CaseStatusMatching, false},
		{"DRAFT to RESOLVED", CaseStatusDraft, CaseStatusResolved, false},
		{"SUBMITTED to MATCHING", CaseStatusSubmitted, CaseStatusMatching, true},
		{"SUBMITTED to CANCELLED", CaseStatusSubmitted, CaseStatusCancelled, false},
		{"MATCHING to MATCHED", CaseStatusMatching, CaseStatusMatched, true},
		{"MATCHING to CANCELLED", CaseStatusMatching, CaseStatusCancelled, true},
		{"MATCHED to IN_PROGRESS", CaseStatusMatched, CaseStatusInProgress, true},
		{"MATCHED to CANCELLED", CaseStatusMatched, CaseStatusCancelled, false},
		{"IN_PROGRESS to RESOLVED", CaseStatusInProgress, CaseStatusResolved, true},
		{"IN_PROGRESS to CANCELLED", CaseStatusInProgress, CaseStatusCancelled, true},
		{"RESOLVED to DRAFT", CaseStatusResolved, CaseStatusDraft, false},
		{"CANCELLED to DRAFT", CaseStatusCancelled, CaseStatusDraft, false},
		{"Unknown status", CaseStatus("UNKNOWN"), CaseStatusDraft, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			can := tt.current.CanTransitionTo(tt.next)
			require.Equal(t, tt.expectCan, can)
		})
	}
}
