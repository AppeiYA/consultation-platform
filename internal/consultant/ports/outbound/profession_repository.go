package outbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

type ProfessionRepository interface {
	GetProfessionByID(ctx context.Context, professionID string) (*domain.Profession, error)
	GetAllProfessions(ctx context.Context) ([]*domain.Profession, error)
}