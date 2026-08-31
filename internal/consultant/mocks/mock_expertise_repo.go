package mocks

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
)

type MockExpertiseRepository struct {
	SaveManyFn           func(ctx context.Context, expertises []*domain.Expertise) error
	AddFn                func(ctx context.Context, expertise *domain.Expertise) error
	FindByConsultantIDFn func(ctx context.Context, consultantID string) ([]*domain.Expertise, error)
	DeleteFn             func(ctx context.Context, consultantID string, expertiseID string) error
	ReplaceAllFn         func(ctx context.Context, consultantID string, expertises []*domain.Expertise) error
}

func (m *MockExpertiseRepository) SaveMany(ctx context.Context, expertises []*domain.Expertise) error {
	if m.SaveManyFn != nil {
		return m.SaveManyFn(ctx, expertises)
	}
	return nil
}

func (m *MockExpertiseRepository) Add(ctx context.Context, expertise *domain.Expertise) error {
	if m.AddFn != nil {
		return m.AddFn(ctx, expertise)
	}
	return nil
}

func (m *MockExpertiseRepository) FindByConsultantID(ctx context.Context, consultantID string) ([]*domain.Expertise, error) {
	if m.FindByConsultantIDFn != nil {
		return m.FindByConsultantIDFn(ctx, consultantID)
	}
	return nil, nil
}

func (m *MockExpertiseRepository) Delete(ctx context.Context, consultantID string, expertiseID string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, consultantID, expertiseID)
	}
	return nil
}

func (m *MockExpertiseRepository) ReplaceAll(ctx context.Context, consultantID string, expertises []*domain.Expertise) error {
	if m.ReplaceAllFn != nil {
		return m.ReplaceAllFn(ctx, consultantID, expertises)
	}
	return nil
}

var _ outbound.ExpertiseRepository = (*MockExpertiseRepository)(nil)
