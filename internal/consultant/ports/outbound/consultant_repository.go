package outbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

type ConsultantRepository interface {
	Save(ctx context.Context, consultant *domain.Consultant) error
	Update(ctx context.Context, consultant *domain.Consultant) error

	FindByID(ctx context.Context, id string) (*domain.Consultant, error)
	FindByUserID(ctx context.Context, userID string) (*domain.Consultant, error)

	ExistsByUserID(ctx context.Context, userID string) (bool, error)
}