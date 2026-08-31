package mocks

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/domain"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/outbound"
)

type MockMatchingRunRepository struct {
	SaveFn               func(ctx context.Context, run *domain.MatchingRun) error
	FindByIDFn           func(ctx context.Context, id string) (*domain.MatchingRun, error)
	FindLatestByCaseIDFn func(ctx context.Context, caseID string) (*domain.MatchingRun, error)
	FindByCaseIDFn       func(ctx context.Context, caseID string) ([]*domain.MatchingRun, error)
}

func (m *MockMatchingRunRepository) Save(ctx context.Context, run *domain.MatchingRun) error {
	if m.SaveFn != nil {
		return m.SaveFn(ctx, run)
	}
	return nil
}

func (m *MockMatchingRunRepository) FindByID(ctx context.Context, id string) (*domain.MatchingRun, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockMatchingRunRepository) FindLatestByCaseID(ctx context.Context, caseID string) (*domain.MatchingRun, error) {
	if m.FindLatestByCaseIDFn != nil {
		return m.FindLatestByCaseIDFn(ctx, caseID)
	}
	return nil, nil
}

func (m *MockMatchingRunRepository) FindByCaseID(ctx context.Context, caseID string) ([]*domain.MatchingRun, error) {
	if m.FindByCaseIDFn != nil {
		return m.FindByCaseIDFn(ctx, caseID)
	}
	return nil, nil
}

var _ outbound.MatchingRunRepository = (*MockMatchingRunRepository)(nil)
