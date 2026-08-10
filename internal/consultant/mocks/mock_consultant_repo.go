package mocks

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

type MockConsultantRepository struct {
	SaveFn func(ctx context.Context, consultant *domain.Consultant) error
	ExistsByUserIDFn func(ctx context.Context, userID string) (bool, error)
	UpdateFn func(ctx context.Context, consultant *domain.Consultant) error
	FindByIDFn func(ctx context.Context, id string) (*domain.Consultant, error)
	FindByUserIDFn func(ctx context.Context, userID string) (*domain.Consultant, error)
}

func (m *MockConsultantRepository) Save(ctx context.Context, consultant *domain.Consultant) error {
	if m.SaveFn != nil {
		return m.SaveFn(ctx, consultant)
	}
	return nil
}

func (m *MockConsultantRepository) ExistsByUserID(ctx context.Context, userID string) (bool, error) {
	if m.ExistsByUserIDFn != nil {
		return m.ExistsByUserIDFn(ctx, userID)
	}
	return false, nil
}

func (m *MockConsultantRepository) Update(ctx context.Context, consultant *domain.Consultant) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, consultant)
	}
	return nil
}

func (m *MockConsultantRepository) FindByID(ctx context.Context, id string) (*domain.Consultant, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockConsultantRepository) FindByUserID(ctx context.Context, userID string) (*domain.Consultant, error) {
	if m.FindByUserIDFn != nil {
		return m.FindByUserIDFn(ctx, userID)
	}
	return nil, nil
}