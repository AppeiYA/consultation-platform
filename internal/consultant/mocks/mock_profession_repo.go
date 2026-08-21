package mocks

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

type MockProfessionRepository struct {
	GetProfessionByIDFn func(ctx context.Context, professionID string) (*domain.Profession, error)
	GetAllProfessionsFn func(ctx context.Context) ([]*domain.Profession, error)
}

func (m *MockProfessionRepository) GetProfessionByID(ctx context.Context, professionID string) (*domain.Profession, error) {
	if m.GetProfessionByIDFn != nil {
		return m.GetProfessionByIDFn(ctx, professionID)
	}
	return nil, nil
}

func (m *MockProfessionRepository) GetAllProfessions(ctx context.Context) ([]*domain.Profession, error) {
	if m.GetAllProfessionsFn != nil {
		return m.GetAllProfessionsFn(ctx)
	}
	return nil, nil
}