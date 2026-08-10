package inbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

type GetConsultantInt interface {
	ByID(ctx context.Context, id string) (*domain.Consultant, error)
	ByUserID(ctx context.Context, userID string) (*domain.Consultant, error)
}