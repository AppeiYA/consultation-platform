package outbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/identity/domain"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error

	Update(ctx context.Context, user *domain.User) error

	FindByID(ctx context.Context, id string) (*domain.User, error)

	FindByEmail(ctx context.Context, email domain.Email) (*domain.User, error)

	Delete(ctx context.Context, id string) error
}