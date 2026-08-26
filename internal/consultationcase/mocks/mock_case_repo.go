package mocks

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultationcase/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultationcase/ports/outbound"
)

type MockCaseRepository struct {
	SaveCaseFn            func(ctx context.Context, c *domain.ConsultationCase) error
	FindCaseByIDFn        func(ctx context.Context, id string) (*domain.ConsultationCase, error)
	FindCasesByClientIDFn func(ctx context.Context, clientID string) ([]*domain.ConsultationCase, error)
	UpdateCaseFn          func(ctx context.Context, c *domain.ConsultationCase) error
	DeleteCaseFn          func(ctx context.Context, id string) error
}

func (m *MockCaseRepository) SaveCase(ctx context.Context, c *domain.ConsultationCase) error {
	if m.SaveCaseFn != nil {
		return m.SaveCaseFn(ctx, c)
	}
	return nil
}

func (m *MockCaseRepository) FindCaseByID(ctx context.Context, id string) (*domain.ConsultationCase, error) {
	if m.FindCaseByIDFn != nil {
		return m.FindCaseByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockCaseRepository) FindCasesByClientID(ctx context.Context, clientID string) ([]*domain.ConsultationCase, error) {
	if m.FindCasesByClientIDFn != nil {
		return m.FindCasesByClientIDFn(ctx, clientID)
	}
	return nil, nil
}

func (m *MockCaseRepository) UpdateCase(ctx context.Context, c *domain.ConsultationCase) error {
	if m.UpdateCaseFn != nil {
		return m.UpdateCaseFn(ctx, c)
	}
	return nil
}

func (m *MockCaseRepository) DeleteCase(ctx context.Context, id string) error {
	if m.DeleteCaseFn != nil {
		return m.DeleteCaseFn(ctx, id)
	}
	return nil
}

var _ outbound.CaseRepository = (*MockCaseRepository)(nil)
