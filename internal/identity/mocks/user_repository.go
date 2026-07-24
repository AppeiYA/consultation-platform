package mocks

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/identity/domain"
)

type MockUserRepository struct {
	SaveFn        func(ctx context.Context, user *domain.User) error
	UpdateFn      func(ctx context.Context, user *domain.User) error
	FindByIDFn    func(ctx context.Context, id string) (*domain.User, error)
	FindByEmailFn func(ctx context.Context, email domain.Email) (*domain.User, error)
	DeleteFn      func(ctx context.Context, id string) error
}

func (m *MockUserRepository) Save(ctx context.Context, user *domain.User) error {
	return m.SaveFn(ctx, user)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	return m.UpdateFn(ctx, user)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	return m.FindByIDFn(ctx, id)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email domain.Email) (*domain.User, error) {
	return m.FindByEmailFn(ctx, email)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	return m.DeleteFn(ctx, id)
}
