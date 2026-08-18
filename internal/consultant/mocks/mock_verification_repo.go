package mocks

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

type MockVerificationRepository struct {
	SaveFn func(ctx context.Context, verification *domain.ConsultantVerification) error
	FindByIDFn func(ctx context.Context, id string) (*domain.ConsultantVerification, error)

	FindByConsultantIDFn func(ctx context.Context, consultantID string) (*domain.ConsultantVerification, error)

	UpdateFn func(ctx context.Context, verification *domain.ConsultantVerification) error
}

func (m *MockVerificationRepository) Save(ctx context.Context, verification *domain.ConsultantVerification) error {
	if m.SaveFn != nil {
		return m.SaveFn(ctx, verification)
	}
	return nil
}

func (m *MockVerificationRepository) FindByID(ctx context.Context, id string) (*domain.ConsultantVerification, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockVerificationRepository) FindByConsultantID(ctx context.Context, consultantID string) (*domain.ConsultantVerification, error) {
	if m.FindByConsultantIDFn != nil {
		return m.FindByConsultantIDFn(ctx, consultantID)
	}
	return nil, nil
}

func (m *MockVerificationRepository) Update(ctx context.Context, verification *domain.ConsultantVerification) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, verification)
	}
	return nil
}