package inbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

type ListProfessionsInt interface {
	Execute(ctx context.Context) ([]*domain.Profession, error)
}