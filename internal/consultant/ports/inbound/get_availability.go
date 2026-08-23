package inbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
)

type GetAvailabilityInt interface {
	Execute(ctx context.Context, consultantID string) ([]*domain.ConsultantAvailability, error)
}