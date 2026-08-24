package inbound

import (
	"context"

	"github.com/AppeiYA/consultation-platform/internal/consultant/domain"
	"github.com/AppeiYA/consultation-platform/internal/consultant/usecase/dto"
)

type UpdateAvailabilityInt interface {
	Execute(ctx context.Context, userID string, req *dto.UpdateAvailabilityRequest) (*domain.ConsultantAvailability, error)
}