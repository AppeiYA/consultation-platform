package outbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
)

type CaseDetails struct {
	CaseID      string
	ClientID    string
	Category    domain.MatchingCategory
	Title       string
	Description string
}

type CaseReader interface {
	GetCaseDetails(ctx context.Context, caseID string) (*CaseDetails, error)
}
