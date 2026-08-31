package mocks

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/outbound"
)

type MockCaseReader struct {
	GetCaseDetailsFn func(ctx context.Context, caseID string) (*outbound.CaseDetails, error)
}

func (m *MockCaseReader) GetCaseDetails(ctx context.Context, caseID string) (*outbound.CaseDetails, error) {
	if m.GetCaseDetailsFn != nil {
		return m.GetCaseDetailsFn(ctx, caseID)
	}
	return nil, nil
}

var _ outbound.CaseReader = (*MockCaseReader)(nil)
