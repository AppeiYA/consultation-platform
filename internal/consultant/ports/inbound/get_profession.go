package inbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

type GetProfessionInt interface {
	Execute(ctx context.Context, professionID string) (*domain.Profession, error)
}